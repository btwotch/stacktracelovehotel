package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s [filename]\n", os.Args[0])
	}

	fmt.Printf("%s", dotPreamble("graph"))

	for _, filename := range os.Args[1:] {

		b, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}

		st := parseStacktrace(b)
		st.name = filepath.Base(filename)

		//		st.ToDot()
		fmt.Printf("%s", st.dotRelations())
	}

	fmt.Printf("%s", dotClosing())
}

func parseStacktrace(b []byte) stacktrace {
	var err error

	goroutineRegex := regexp.MustCompile("^goroutine ([0-9]+).*:$")
	functionRegex := regexp.MustCompile(`(.*)\.(.*)\(.*\)`)
	lineRegex := regexp.MustCompile(`\s*(\/.*\.go):(\d+)`)

	var st stacktrace
	var fn *function
	for _, lineBytes := range bytes.Split(b, []byte("\n")) {
		line := string(lineBytes)
		match := goroutineRegex.FindStringSubmatch(line)
		if len(match) == 2 {
			st.goroutine, err = strconv.ParseUint(match[1], 10, 64)
			if err != nil {
				panic(err)
			}
			continue
		}
		match = functionRegex.FindStringSubmatch(line)
		if len(match) == 3 {
			fn = &function{}
			fn.packageName = match[1]
			fn.functionName = match[2]
			continue
		}
		match = lineRegex.FindStringSubmatch(line)
		if len(match) == 3 {
			fn.path = match[1]
			fn.line, err = strconv.ParseUint(match[2], 10, 64)
			if err != nil {
				panic(err)
			}
			st.appendFunction(*fn)
			continue
		}
	}
	return st
}
