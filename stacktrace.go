package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type stacktrace struct {
	name      string
	goroutine uint64
	fs        []function
}

func (st *stacktrace) Name() string {
	if st.name != "" {
		return st.name
	}

	return fmt.Sprintf("goroutine %d", st.goroutine)
}

func (st *stacktrace) String() string {
	ret := fmt.Sprintf("goroutine: %d\n", st.goroutine)
	for _, fn := range st.fs {
		ret += fmt.Sprintf("%s: %s.%s:%d\n", fn.path, fn.packageName, fn.functionName, fn.line)
	}

	return ret
}

func (st *stacktrace) ToDot() {
	preamble := dotPreamble(st.Name())

	fmt.Printf("%s", preamble)

	relations := st.dotRelations()

	fmt.Printf("%s", relations)

	closing := dotClosing()
	fmt.Printf("%s", closing)
}

func dotClosing() string {
	closing := fmt.Sprintf("}\n")
	return closing
}

func dotPreamble(title string) string {
	var preamble string
	preamble = fmt.Sprintf("digraph \"%s\"\n", title)
	preamble += fmt.Sprintf("{\n")
	return preamble
}

type relation struct {
	from string
	to   string
}

type relations struct {
	rs map[relation]struct{}
}

func newRelations() relations {
	var rs relations

	rs.rs = make(map[relation]struct{})

	return rs
}
func (rs relations) String() string {
	var ret string
	for relation := range rs.rs {
		from := relation.from
		to := relation.to
		ret += fmt.Sprintf("\t\"%s\" -> \"%s\";\n", from, to)
	}

	return ret
}

func (rs *relations) ingest(other relations) {
	for relation := range other.rs {
		rs.rs[relation] = struct{}{}
	}
}

func (st *stacktrace) dotRelations() relations {
	rs := newRelations()
	for i := 1; i < len(st.fs); i++ {
		from := st.fs[i].String()
		to := st.fs[i-1].String()
		from = strings.TrimSpace(from)
		to = strings.TrimSpace(to)

		rl := relation{
			from: from,
			to:   to,
		}
		rs.rs[rl] = struct{}{}
	}
	return rs
}

func (st *stacktrace) appendFunction(fn function) {
	st.fs = append(st.fs, fn)
}

func parseStacktrace(b []byte) stacktrace {
	var err error

	goroutineRegex := regexp.MustCompile("^goroutine ([0-9]+).*:$")
	functionRegex := regexp.MustCompile(`([^/]*)\.([^/]*)\(.*\)`)
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
			if !fn.Empty() {
				st.appendFunction(*fn)
			}
			fn = &function{}
			continue
		}

	}

	return st
}
