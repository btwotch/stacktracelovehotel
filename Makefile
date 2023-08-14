st: *.go
	go fmt
	goimports -w .
	go build -o st
