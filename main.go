package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s [filename]\n", os.Args[0])
	}

	var sts stacktraces

	for _, filename := range os.Args[1:] {

		b, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}

		name := filepath.Base(filename)
		sts.addFromString(name, b)

	}

	fmt.Printf("%s\n", sts.String())
}
