package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/yookoala/realpath"
)

const (
	helpString = `
usage: overalls -project=[path] -covermode[mode] OPTIONS

overalls recursively traverses your projects directory structure
running 'go test -covermode=count -coverprofile=profile.coverprofile'
in each directory with go test files, concatenates them into one
coverprofile in your root directory named 'overalls.coverprofile'

OPTIONS
  -project
	Your project path relative to the '$GOPATH/src' directory
	example: -project=github.com/go-playground/overalls

  -covermode
    Mode to run when testing files.
    default:count

OPTIONAL

  -go-binary
    An alternative 'go' binary to run the tests, for example to use 'richgo' for
		more human-friendly output.
    example: -go-binary=richgo
    default: 'go'

  -ignore
    A comma separated list of directory names to ignore, relative to project path.
    example: -ignore=[.git,.hiddentdir...]
    default: '.git,vendor'

  -debug
    A flag indicating whether to print debug messages.
    example: -debug
    default:false

  -concurrency
    Limit the number of packages being processed at one time.
    The minimum value must be 2 or more when set.
    example: -concurrency=5
    default: unlimited
`
)

const (
	defaultIgnores = ".git,vendor"
	outFilename    = "overalls.coverprofile"
	pkgFilename    = "profile.coverprofile"
)

var (
	modeRegex       = regexp.MustCompile("mode: [a-z]+\n")
	srcPath         string
	projectPath     string
	goBinary        string
	ignoreFlag      string
	projectFlag     string
	coverFlag       string
	helpFlag        bool
	debugFlag       bool
	concurrencyFlag int
	isLimited       bool
	emptyStruct     struct{}
	ignores         = map[string]struct{}{}
	flagArgs        []string
)

func help() {
	fmt.Printf(helpString)
}

func init() {
	flag.StringVar(&goBinary, "go-binary", "go", "Use an alternative test runner such as 'richgo'")
	flag.StringVar(&projectFlag, "project", "", "-project [path]: relative to the '$GOPATH/src' directory")
	flag.StringVar(&coverFlag, "covermode", "count", "Mode to run when testing files")
	flag.StringVar(&ignoreFlag, "ignore", defaultIgnores, "-ignore [dir1,dir2...]: comma separated list of directory names to ignore")
	flag.IntVar(&concurrencyFlag, "concurrency", -1, "-concurrency [int]: number of packages to process concurrently, The minimum value must be 2 or more when set.")
	flag.BoolVar(&debugFlag, "debug", false, "-debug [true|false]")
	flag.BoolVar(&helpFlag, "help", false, "-help")
}

func parseFlags(logger *log.Logger) {

	flag.Parse()

	if helpFlag {
		help()
		os.Exit(0)
	}

	if debugFlag {
		fmt.Println("GOPATH:", os.Getenv("GOPATH"))
	}

	fmt.Println("|", projectFlag)
	projectFlag = filepath.Clean(projectFlag)

	if debugFlag {
		fmt.Println("Project Path:", projectFlag)
	}

	if len(projectFlag) == 0 || projectFlag == "." {
		fmt.Printf("\n**invalid project path '%s'\n", projectFlag)
		help()
		os.Exit(1)
	}

	pkg, err := build.Default.Import(projectFlag, "", build.FindOnly)
	if err != nil {
		fmt.Printf("\n**could not find project path '%s' in GOPATH '%s'\n", projectFlag, os.Getenv("GOPATH"))
		os.Exit(1)
	}
	srcPath = pkg.SrcRoot

	flagArgs = flag.Args()

	switch coverFlag {
	case "set", "atomic":
	case "count":
		for _, flg := range flagArgs {
			if flg == "-race" {
				logger.Println("\n*****\n** WARNING: some common patterns in parallel code can trigger race conditions when using coverprofile=count and the -race flag; in which case coverprofile=atomic should be used.\n*****")
				break
			}
		}
	default:
		fmt.Printf("\n**invalid covermode '%s'\n", coverFlag)
		os.Exit(1)
	}

	arr := strings.Split(ignoreFlag, ",")
	for _, v := range arr {
		ignores[v] = emptyStruct
	}

	isLimited = concurrencyFlag != -1

	if isLimited && concurrencyFlag < 1 {
		fmt.Printf("\n**invalid concurrency value '%d', value must be at least 1\n", concurrencyFlag)
		os.Exit(1)
	}
}

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	runMain(logger)
}

func runMain(logger *log.Logger) {
	parseFlags(logger)

	var err error
	var wd string

	projectPath = filepath.Join(srcPath, projectFlag)

	if err = os.Chdir(projectPath); err != nil {
		logger.Printf("\n**invalid project path '%s'\n%s\n", projectFlag, err)
		help()
		os.Exit(1)
	}

	if debugFlag {
		wd, err = os.Getwd()
		if err != nil {
			fmt.Println(err)
		}

		logger.Println("Working DIR:", wd)
	}

	testFiles(logger)
}

func scanOutput(r io.ReadCloser, fn func(...interface{})) {
	defer r.Close()
	bs := bufio.NewScanner(r)
	for bs.Scan() {
		fn(bs.Text())
	}
	if err := bs.Err(); err != nil {
		fn(fmt.Sprintf("Scan error: %v", err.Error()))
	}
}

func processDIR(logger *log.Logger, wg *sync.WaitGroup, fullPath, relPath string, out chan<- []byte, semaphore chan struct{}) {
	defer wg.Done()

	if isLimited {
		semaphore <- struct{}{}
	}

	// 1 for "test", 4 for covermode, coverprofile, outputdir, relpath
	args := make([]string, 1, 1+len(flagArgs)+4)
	args[0] = "test"
	// To split '-- <go test arguments> -args <program arguments>'
	for i, arg := range flagArgs {
		if arg == "-args" {
			args = append(args, flagArgs[:i]...)
			flagArgs = flagArgs[i:]
			break
		}
	}
	args = append(args, "-covermode="+coverFlag, "-coverprofile="+pkgFilename, "-outputdir="+fullPath+"/", relPath)
	args = append(args, flagArgs...)
	fmt.Printf("Test args: %+v\n", args)

	cmd := exec.Command(goBinary, args...)
	if debugFlag {
		logger.Println("Processing:", strings.Join(cmd.Args, " "))
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Fatal("Unable to get process stdout")
	}

	go scanOutput(stdout, logger.Print)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Fatal("Unable to get process stderr")
	}

	go scanOutput(stderr, logger.Print)

	if err := cmd.Run(); err != nil {
		logger.Fatal("ERROR:", err)
	}

	b, err := ioutil.ReadFile(relPath + "/profile.coverprofile")
	if err != nil {
		logger.Fatal("ERROR:", err)
	}

	out <- b

	if isLimited {
		<-semaphore
	}
}

// walk is like filepath.Walk, but it follows symlinks and only calls walkFunc on directories.
func walkDirectories(path string, walkFunc func(path string, info os.FileInfo) error) error {
	seen := make(map[string]bool)

	var walkHelper func(path string) error
	walkHelper = func(path string) error {
		qualifiedPath, err := realpath.Realpath(path)
		if err != nil {
			return err
		}

		// Prevent circular links.
		if seen[qualifiedPath] {
			return nil
		}
		seen[qualifiedPath] = true

		// Skip anything that isn't a directory.
		file, err := os.Stat(path)
		if err != nil {
			return err
		}
		if !file.IsDir() {
			return nil
		}

		err = walkFunc(path, file)
		if err != nil {
			if err == filepath.SkipDir {
				return nil
			}
			return err
		}

		files, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}

		for _, file := range files {
			err = walkHelper(filepath.Join(path, file.Name()))
			if err != nil {
				return err
			}
		}
		return nil
	}
	return walkHelper(path)
}

func testFiles(logger *log.Logger) {

	var semaphore chan struct{}

	if isLimited {
		semaphore = make(chan struct{}, concurrencyFlag)
	}

	out := make(chan []byte)
	wg := &sync.WaitGroup{}

	walker := func(path string, info os.FileInfo) error {
		rel, err := filepath.Rel(projectPath, path)
		if err != nil {
			logger.Fatalf("Could not make path '%s' relative to project path '%s'", path, projectPath)
		}

		if _, ignore := ignores[rel]; ignore {
			return filepath.SkipDir
		}

		rel = "./" + rel

		if files, err := filepath.Glob(rel + "/*_test.go"); len(files) == 0 || err != nil {

			if err != nil {
				logger.Fatal("Error checking for test files")
			}

			if debugFlag {
				logger.Println("No Go Test files in DIR:", rel, "skipping")
			}

			return nil
		}

		wg.Add(1)

		go processDIR(logger, wg, path, rel, out, semaphore)

		return nil
	}

	if err := walkDirectories(projectPath, walker); err != nil {
		logger.Fatalf("\n**could not walk project path '%s'\n%s\n", projectPath, err)
	}

	go func() {
		wg.Wait()
		close(out)

		if isLimited {
			close(semaphore)
		}
	}()

	buff := bytes.NewBufferString("")

	for cover := range out {
		buff.Write(cover)
	}

	final := buff.String()
	final = modeRegex.ReplaceAllString(final, "")
	final = "mode: " + coverFlag + "\n" + final

	if err := ioutil.WriteFile(outFilename, []byte(final), 0644); err != nil {
		logger.Fatal("ERROR Writing \""+outFilename+"\"", err)
	}
}
