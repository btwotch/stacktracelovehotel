package main

import "fmt"

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

func (st *stacktrace) dotRelations() string {
	var relations string
	for i := 1; i < len(st.fs); i++ {
		from := st.fs[i].String()
		to := st.fs[i-1].String()
		from = strings.TrimSpace(from)
		to = strings.TrimSpace(to)
		relations += fmt.Sprintf("\t\"%s\" -> \"%s\";\n", from, to)
	}
	return relations
}

func (st *stacktrace) appendFunction(fn function) {
	st.fs = append(st.fs, fn)
}
