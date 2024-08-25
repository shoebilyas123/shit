package main

import (
	"fmt"
	"os"

	"github.com/shoebilyas123/shit/commands"
	"github.com/shoebilyas123/shit/initialize"
	"github.com/shoebilyas123/shit/plumbings"
)

// The init pattern that we are looking for: shit (init)(argvs[1])
func main() {
	argvs := os.Args

	if len(argvs) < 2 {
		fmt.Println("Shit is a content-addressable file system inspired by git")
		os.Exit(0)
	}

	base_cmd := argvs[1]
	switch base_cmd {
	case commands.INIT:
		initialize.Init(argvs[2:])
	case commands.HASH:
		plumbings.HashObject(argvs[2:], "blob", "")
	case commands.UPDATE_INDEX:
		plumbings.UpdateIndex(argvs[2:])
	case commands.WRITE_TREE:
		plumbings.WriteTree()
	case commands.CAT_FILE:
		plumbings.CatFile(argvs[2:])
	case commands.COMMIT_TREE:
		if len(argvs[2:]) == 0 {
			fmt.Println("Fatal: Empty commit-tree command")
			return
		}
		tree_sha1 := argvs[2]
		parent_sha1 := ""
		commit_msg := ""

		if len(argvs[3:]) < 2 {
			fmt.Println("Fatal: Invalid commit-tree command")
			return
		}
		if argvs[3] == "-p" {
			if argvs[4] == "-m" || len(argvs[4]) <= 0 {
				fmt.Println("fatal: parent sha_1 not provided")
				return
			}
			parent_sha1 = argvs[4]
		} else if argvs[3] == "-m" {
			if len(argvs[4]) <= 0 {
				fmt.Println("Fatal: Commit message is required")
				return
			}
			commit_msg = argvs[4]
		} else {

			if argvs[3] == "-p" && argvs[5] != "-m" {
				fmt.Println("Fatal: Commit message required")
				return
			} else {
				commit_msg = argvs[6]
			}
		}

		plumbings.CommitTree(tree_sha1, parent_sha1, commit_msg)
	}
}
