package main

import (
	"fmt"
	"os"
)

var (
	y = 2
)

func main() {
	fmt.Fprintln(os.Stdout, "Hello")
}

func a(x int) int {
	return x * y
}
