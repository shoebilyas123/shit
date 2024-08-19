package main

import (
	"fmt"
	"os"

	"github.com/shoebilyas123/shit/commands"
	"github.com/shoebilyas123/shit/initialize"
)

// The init pattern that we are looking for: shit (init)(argvs[1])
func main() {
	argvs := os.Args

	if len(argvs) < 2 {
		fmt.Println("Shit is a content-addressable file system inspired by git")
		os.Exit(0)
	}

	base_cmd := argvs[1]
	fmt.Printf("%+v\n", argvs[2:])
	switch base_cmd {
	case commands.INIT:
		initialize.Init(argvs[2:])
	}
}