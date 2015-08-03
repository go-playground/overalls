Package overalls
================

[![Build Status](https://semaphoreci.com/api/v1/projects/a0634385-0174-48b1-9190-fb2645fdb9a2/501046/badge.svg)](https://semaphoreci.com/joeybloggs/overalls)
[![GoDoc](https://godoc.org/github.com/bluesuncorp/overalls?status.svg)](https://godoc.org/github.com/bluesuncorp/overalls)

Package overalls takes multi-package go projects, runs test coverage tests on all packages in each directory and finally concatenates into a single file for tools like goveralls.

Usage and documentation
------

```go
//Example 
overalls -project=github.com/bluesuncorp/overalls" -covermode=count -debug
```
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
    
How to Contribute
------

There will always be a development branch. In order to contribute, 
please make your pull requests against that branch.

If the changes being proposed or requested are breaking changes, please create an issue.

License
------
Distributed under MIT License, please see license file in code for more details.
