package main

import "fmt"

type function struct {
	path         string
	line         uint64
	packageName  string
	functionName string
}

func (fn *function) String() string {
	return fmt.Sprintf("%s.%s", fn.packageName, fn.functionName)
}

func (fn *function) Empty() bool {
	return fn.functionName == ""
}
