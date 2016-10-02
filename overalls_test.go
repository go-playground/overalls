package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	. "gopkg.in/go-playground/assert.v1"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called
//
// or
//
// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html
//

func TestOveralls_Default(t *testing.T) {
	withTestingOveralls(t, func(output []byte, fileBytes []byte) {
		final := string(fileBytes)
		NotEqual(t, strings.Index(final, "main.go"), -1)
		NotEqual(t, strings.Index(final, "test-files/good/main.go"), -1)
		NotEqual(t, strings.Index(final, "test-files/good2/main.go"), -1)
		MatchRegex(t, string(output), "-covermode=count")
		MatchRegex(t, string(output), "-outputdir=.*/go-playground/overalls/test-files/good")

		MatchRegex(t, string(output), "go test -covermode=count")
	})
}

func TestOveralls_WithExtraArguments(t *testing.T) {
	withTestingOveralls(t, func(output []byte, fileBytes []byte) {
		MatchRegex(t, string(output), "Processing: go test")

		MatchRegex(t, string(output), "=== RUN")
		MatchRegex(t, string(output), "--- PASS: TestGood")
	}, "--", "-v")
}

func withTestingOveralls(t *testing.T, fn func(output []byte, coverage []byte), args ...string) {
	baseArgs := []string{"-project=github.com/go-playground/overalls/test-files", "-covermode=count", "-debug"}
	args = append(baseArgs, args...)
	args = append([]string{"overalls"}, args...)
	defer stubArgs(args...)()

	out := &bytes.Buffer{}
	runMain(log.New(out, "", 0))

	fileBytes, err := ioutil.ReadFile(srcPath + "github.com/go-playground/overalls/test-files/overalls.coverprofile")

	// Ensure no error reading file
	Equal(t, err, nil)
	fn(out.Bytes(), fileBytes)
}

func stubArgs(args ...string) func() {
	oldArgs := os.Args
	os.Args = args
	return func() {
		os.Args = oldArgs
	}
}
