/*
Package overalls takes multi-package go projects, runs test coverage tests on
all packages in each directory and finally concatenates into a single file for
tools like goveralls.

	$ overalls -help

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

	  -concurrency
	    Limit the number of packages being processed at one time.
	    The minimum value must be 2 or more when set.
	    example: -concurrency=5
	    default: unlimited
*/
package main
