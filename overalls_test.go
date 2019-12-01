package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
		NotEqual(t, strings.Index(final, "testdata/circular/lib/main.go"), -1)
		NotEqual(t, strings.Index(final, "testdata/good/main.go"), -1)
		NotEqual(t, strings.Index(final, "testdata/good2/main.go"), -1)
		NotEqual(t, strings.Index(final, "testdata/symlink-real-folder/main.go"), -1)
		MatchRegex(t, string(output), "-covermode=atomic")
		MatchRegex(t, string(output), "-outputdir=.*/go-playground/overalls/testdata/circular/lib")
		MatchRegex(t, string(output), "-outputdir=.*/go-playground/overalls/testdata/good")
		MatchRegex(t, string(output), "-outputdir=.*/go-playground/overalls/testdata/good2")
		MatchRegex(t, string(output), "-outputdir=.*/go-playground/overalls/testdata/symlink-real-folder")
		MatchRegex(t, string(output), "go test -covermode=atomic")
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

	baseArgs := []string{"-project=github.com/go-playground/overalls/testdata", "-covermode=atomic", "-debug"}
	args = append(baseArgs, args...)
	args = append([]string{"overalls"}, args...)
	defer stubArgs(args...)()

	out := &buffer{}
	runMain(log.New(out, "", 0))

	fileBytes, err := ioutil.ReadFile(filepath.Join(projectPath, "overalls.coverprofile"))

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

type buffer struct {
	b bytes.Buffer
	m sync.Mutex
}

func (b *buffer) Read(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Read(p)
}
func (b *buffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Write(p)
}
func (b *buffer) String() string {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.String()
}

func (b *buffer) Bytes() []byte {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Bytes()
}
func (b *buffer) Cap() int {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Cap()
}
func (b *buffer) Grow(n int) {
	b.m.Lock()
	defer b.m.Unlock()
	b.b.Grow(n)
}
func (b *buffer) Len() int {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Len()
}
func (b *buffer) Next(n int) []byte {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Next(n)
}
func (b *buffer) ReadByte() (c byte, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.ReadByte()
}
func (b *buffer) ReadBytes(delim byte) (line []byte, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.ReadBytes(delim)
}
func (b *buffer) ReadFrom(r io.Reader) (n int64, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.ReadFrom(r)
}
func (b *buffer) ReadRune() (r rune, size int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.ReadRune()
}
func (b *buffer) ReadString(delim byte) (line string, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.ReadString(delim)
}
func (b *buffer) Reset() {
	b.m.Lock()
	defer b.m.Unlock()
	b.b.Reset()
}
func (b *buffer) Truncate(n int) {
	b.m.Lock()
	defer b.m.Unlock()
	b.b.Truncate(n)
}
func (b *buffer) UnreadByte() error {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.UnreadByte()
}
func (b *buffer) UnreadRune() error {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.UnreadRune()
}
func (b *buffer) WriteByte(c byte) error {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.WriteByte(c)
}
func (b *buffer) WriteRune(r rune) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.WriteRune(r)
}
func (b *buffer) WriteString(s string) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.WriteString(s)
}
func (b *buffer) WriteTo(w io.Writer) (n int64, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.WriteTo(w)
}
