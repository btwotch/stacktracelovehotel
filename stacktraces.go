package main

type stacktraces struct {
	sts []stacktrace
}

func (sts *stacktraces) String() string {
	str := dotPreamble("Stacktraces")

	rs := newRelations()

	for _, st := range sts.sts {
		rs.ingest(st.dotRelations())
	}

	str += rs.String()

	str += dotClosing()

	return str
}

func (sts *stacktraces) Render(path string) {
	dot := sts.String()
	renderDot(dot, path)
}

func (sts *stacktraces) addFromString(name string, b []byte) {
	st := parseStacktrace(b)
	st.name = name

	sts.sts = append(sts.sts, st)
}
