package main

import (
	"fmt"
	"os"
)

// The init pattern that we are looking for: shit (init)(argvs[1])
func main() {
	argvs := os.Args

	if len(argvs) < 2 {
		fmt.Println("Shit is a content-addressable file system inspired by git")
		os.Exit(0)
	}

}
