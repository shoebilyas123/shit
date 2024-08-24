package index

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/shoebilyas123/shit/common"
	"github.com/shoebilyas123/shit/object"
)

type IndexEntry struct {
	Type     string
	Mode     int64
	SHA_1    string
	Filename string
}

func checkValidMode(mode int64) bool {
	switch mode {
	case 100644:
		fallthrough
	case 040000:
		return true
	default:
		return false
	}
}

func getObjectType(mode int64) string {
	switch mode {
	case 100644:
		return "blob"
	case 040000:
		return "tree"
	default:
		return ""
	}
}

// Scenario 1: This is the first update-index and the index file is empty;
// index entry format: <type: blob|tree> <mode: 100644|040000> <sha_1> <file_name>\n
// we are assuming that the updateindex need to add only the blob
// case 2: we pass a folder to update object to add to the index
func UpdateIndexAdd(argvs []string) {
	var filename, sha_1 string
	var mode int64
	// argvs[1]: [ --fromcache,mode, sha_1, filename ]

	switch argvs[0] {
	case "--fromcache":
		cli_mode, err := strconv.Atoi(argvs[1])
		if err != nil {
			fmt.Printf("fatal: mode must be an integer")
			return
		}

		mode = int64(cli_mode)

		if !checkValidMode(mode) {
			fmt.Printf("fatal: invalid object mode")
			return
		}
		sha_1 = argvs[2]
		filename = argvs[3]
	default:
		filename = argvs[0]
		sha_1 = object.HashObject([]string{filename}, "blob")
		mode = 100644
	}

	if !common.CheckDirExistence("./.shit/index") {
		common.HandleCreateFile("./.shit", "index")
	}

	// Get the index file contents
	// Parse them into an array of []IndexEntry
	// Search for a file with the same file name and update the entire entry
	// Format the new array into a string and read it into the file again
	index_contents, err := os.ReadFile("./.shit/index")

	if err != nil {
		fmt.Printf("fatal: cannot read the index file\n")
		return
	}

	index_file := string(index_contents)
	current_entry := fmt.Sprintf("blob %d %s %s", mode, sha_1, filename)
	update_content := []string{}
	isCurrentEntryUpdated := false
	for _, _entry := range strings.Split(index_file, "\n") {
		if _entry == "" {
			continue
		}

		if strings.Split(_entry, " ")[3] == strings.Split(current_entry, " ")[3] {
			update_content = append(update_content, current_entry)
			isCurrentEntryUpdated = true
		} else {
			update_content = append(update_content, _entry)
		}
	}

	if !isCurrentEntryUpdated {
		update_content = append(update_content, current_entry)
	}
	fmt.Println(strings.Join(update_content, "\n"))
	// fmt.Println(current_entry)s

	if !common.CheckDirExistence("./.shit/index") {
		common.HandleCreateFile("./.shit", "index")
	}

	os.WriteFile("./.shit/index", []byte(strings.Join(update_content, "\n")), 0644)
}

func UpdateIndex(argvs []string) {
	// argvs[0]: [--add]
	// argvs[1]: [ filename ]
	switch argvs[0] {
	case "--add":
		UpdateIndexAdd(argvs[1:])
	default:
		return
	}
}
