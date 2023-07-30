all: stacktrace.png

st: *.go
	go fmt
	goimports -w .
	go build -o st

stacktrace.dot: st
	./st st_*.txt > stacktrace.dot

stacktrace.png: st stacktrace.dot
	dot -Tpng -o stacktrace.png stacktrace.dot
