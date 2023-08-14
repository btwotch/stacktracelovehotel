package main

import (
	"fmt"
	"os/exec"
)

// dot -Tpng -o stacktrace.png stacktrace.dot
func renderDot(dot string, path string) {
	cmd := exec.Command("dot", "-Tpng", "-o", path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	stdin.Write([]byte(dot))
	stdin.Close()
	err = cmd.Run()
	if err != nil {
		output, _ := cmd.CombinedOutput()
		newerror := fmt.Errorf("dot failed: output is:\n%s\nerror is: %+v", output, err)
		panic(newerror)
	}
}
