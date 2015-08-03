package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
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
	example: -project=github.com/bluesuncorp/overalls

  -covermode
    Mode to run when testing files.
    default:count

OPTIONAL

  -ignore
    A comma separated list of directory names to ignore, relative to project path.
    example: -ignore=[.git,.hiddentdir...]
    default: '.git'

  -debug
    A flag indicating whether to print debug messages.
    example: -debug
    default:false
`
)

const (
	defaultIgnores = ".git"
	outFilename    = "overalls.coverprofile"
)

var (
	modeRegex   = regexp.MustCompile("^mode: [a-z]+\n")
	gopath      = filepath.Clean(os.Getenv("GOPATH"))
	srcPath     = gopath + "/src/"
	projectPath string
	ignoreFlag  string
	projectFlag string
	countFlag   string
	helpFlag    bool
	debugFlag   bool
	emptyStruct struct{}
	ignores     = map[string]struct{}{}
)

func help() {
	fmt.Printf(helpString)
}

func parseFlags() {

	flag.StringVar(&projectFlag, "project", "", "-project [path]: relative to the '$GOPATH/src' directory")
	flag.StringVar(&countFlag, "covermode", "count", "Mode to run when testing files")
	flag.StringVar(&ignoreFlag, "ignore", defaultIgnores, "-ignore [dir1,dir2...]: comma separated list of directory names to ignore")
	flag.BoolVar(&debugFlag, "debug", false, "-debug [true|false]")
	flag.BoolVar(&helpFlag, "help", false, "-help")
	flag.Parse()

	if helpFlag {
		help()
		os.Exit(0)
	}

	if debugFlag {
		fmt.Println("GOPATH:", gopath)
	}

	if len(gopath) == 0 || gopath == "." {
		fmt.Printf("\n**invalid GOPATH '%s'\n", gopath)
		os.Exit(1)
	}

	projectFlag = filepath.Clean(projectFlag)

	if debugFlag {
		fmt.Println("Project Path:", projectFlag)
	}

	if len(projectFlag) == 0 || projectFlag == "." {
		fmt.Printf("\n**invalid project path '%s'\n", projectFlag)
		help()
		os.Exit(1)
	}

	arr := strings.Split(ignoreFlag, ",")
	for _, v := range arr {
		ignores[v] = emptyStruct
	}
}

func main() {

	parseFlags()

	var err error
	var wd string

	projectPath = srcPath + projectFlag + "/"

	if err = os.Chdir(projectPath); err != nil {
		fmt.Printf("\n**invalid project path '%s'\n%s\n", projectFlag, err)
		help()
		os.Exit(1)
	}

	if debugFlag {
		wd, err = os.Getwd()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Working DIR:", wd)
	}

	testFiles()
}

func processDIR(wg *sync.WaitGroup, relPath string, out chan<- []byte) {

	defer wg.Done()

	if debugFlag {
		fmt.Println("Processing: go test -covermode=" + countFlag + " -coverprofile=profile.coverprofile " + relPath)
	}

	// go test -covermode=count -coverprofile=count.out
	cmd := exec.Command("go", "test", "-covermode="+countFlag, "-coverprofile=profile.coverprofile", relPath)
	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	b, err := ioutil.ReadFile(relPath + "/profile.coverprofile")
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	out <- b
}

func testFiles() {

	out := make(chan []byte)
	wg := &sync.WaitGroup{}

	walker := func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			return nil
		}

		rel := strings.Replace(path, projectPath, "", 1)

		if _, ignore := ignores[rel]; ignore {
			return filepath.SkipDir
		}

		rel = "./" + rel

		if files, err := filepath.Glob(rel + "/*_test.go"); len(files) == 0 || err != nil {

			if err != nil {
				fmt.Println("Error checking for test files")
				os.Exit(1)
			}

			if debugFlag {
				fmt.Println("No Go Test files in DIR:", rel, "skipping")
			}

			return nil
		}

		wg.Add(1)
		go processDIR(wg, rel, out)

		return nil
	}

	if err := filepath.Walk(projectPath, walker); err != nil {
		fmt.Printf("\n**could not walk project path '%s'\n%s\n", projectPath, err)
		os.Exit(1)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	buff := bytes.NewBufferString("")

	for cover := range out {
		buff.Write(cover)
	}

	final := buff.String()

	if modeRegex.Match(buff.Bytes()) {
		final = modeRegex.ReplaceAllString(final, "")
	}

	final = "mode: " + countFlag + "\n" + final

	if err := ioutil.WriteFile(outFilename, []byte(final), 0644); err != nil {
		fmt.Println("ERROR Writing \""+outFilename+"\"", err)
		os.Exit(1)
	}
}
