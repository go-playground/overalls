Package overalls
================

[![Build Status](https://travis-ci.org/go-playground/overalls.svg?branch=master)](https://travis-ci.org/go-playground/overalls)
[![GoDoc](https://godoc.org/github.com/go-playground/overalls?status.svg)](https://godoc.org/github.com/go-playground/overalls)

Package overalls takes multi-package go projects, runs test coverage tests on all packages in each directory and finally concatenates into a single file for tools like goveralls.

Usage and documentation
------
##### Example
	overalls -project=github.com/go-playground/overalls -covermode=count -debug

##### then with other tools such as goveralls
	goveralls -coverprofile=overalls.coverprofile -service semaphore -repotoken $COVERALLS_TOKEN

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
