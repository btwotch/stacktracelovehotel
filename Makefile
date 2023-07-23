all: stacktrace.png

st: *.go
	go fmt
	goimports -w .
	go build -o st

stacktrace.png: st stacktrace.txt
	./st stacktrace.txt stacktrace2.txt > stacktrace.dot
	dot -Tpng -o stacktrace.png stacktrace.dot
