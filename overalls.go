package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	helpString = `
usage: overalls -project=[path] OPTIONS

overalls recursively traverses your projects directory structure
running 'go test -covermode=count -coverprofile=profile.coverprofile'
in each directory with go test files, concatenates them into one
coverprofile in your root directory named 'overalls.coverprofile' and
then submits you results to coveralls using the goveralls package
https://github.com/mattn/goveralls

OPTIONS
  -project
	Your project path relative to the '$GOPATH/src' directory
	example: -project=github.com/bluesuncorp/overalls

OPTIONAL

  -ignore
    A comma separated list of directory names to ignore.
    example: -ignore=[.git,.hiddentdir...]
    default: '.git'

  -debug=[true|false]
    A flag indicating whether to print debug messages.
    example: -debug
    default:false
`
)

const (
	defaultIgnores = ".git"
)

var (
	gopath      = filepath.Clean(os.Getenv("GOPATH"))
	srcPath     = gopath + "/src/"
	projectPath string
	ignoreFlag  string
	projectFlag string
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

	projectPath = srcPath + projectFlag

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

func testFiles() {
	walker := func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() {
			return nil
		}

		if _, ignore := ignores[info.Name()]; ignore {
			return filepath.SkipDir
		}

		if debugFlag {
			fmt.Println("PROCESSING PATH:", path)
		}

		return nil
	}

	if err := filepath.Walk(projectPath, walker); err != nil {
		fmt.Printf("\n**could not walk project path '%s'\n%s\n", projectPath, err)
		os.Exit(1)
	}
}
