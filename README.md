Package overalls
================

[![Build Status](https://travis-ci.org/go-playground/overalls.svg?branch=master)](https://travis-ci.org/go-playground/overalls)
[![GoDoc](https://godoc.org/github.com/go-playground/overalls?status.svg)](https://godoc.org/github.com/go-playground/overalls)

Package overalls takes multi-package go projects, runs test coverage tests on all packages in each directory and finally concatenates into a single file for tools like goveralls and codecov.io.

Usage and documentation
------
##### Example
	overalls -project=github.com/go-playground/overalls -covermode=count -debug

##### then with other tools such as [goveralls](https://github.com/mattn/goveralls)
	goveralls -coverprofile=overalls.coverprofile -service semaphore -repotoken $COVERALLS_TOKEN
	
##### or [codecov.io](https://github.com/codecov/example-go)
	mv overalls.coverprofile coverage.txt
	export CODECOV_TOKEN=###
 	bash <(curl -s https://codecov.io/bash)


##### note:
1. goveralls and codecover currently do not calculate coverage the same way as `go tool cover` see [here](https://github.com/mattn/goveralls/issues/103) and [here](https://github.com/codecov/example-go/issues/13).

2. overalls (and go test) by default will not calculate coverage "across" packages. E.g. if a test in package A covers code in package B overalls will not count it. You may or may not want this depending on whether you're more concerned about unit test coverage or integration test coverage. To enable add the coverpkg flag.
    `overalls -project=github.com/go-playground/overalls -covermode=count -debug -- -coverpkg=./...`

```shell
$ overalls -help

usage: overalls -project=[path] -covermode[mode] OPTIONS -- TESTOPTIONS

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
```

TESTOPTIONS

  Any flags after `--` will be passed as-is to `go test`.
  For example:

```bash
overalls -project=$PROJECT -debug -- -race -v
```

Will call `go test -race -v` under the hood in addition to the `-coverprofile`
commands.

How to Contribute
------

Make a pull request.

If the changes being proposed or requested are breaking changes, please create an issue.

License
------
Distributed under MIT License, please see license file in code for more details.
