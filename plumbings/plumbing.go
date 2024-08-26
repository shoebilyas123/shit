package plumbings

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shoebilyas123/shit/common"
)

const (
	BLOB_MODE = 100644
	TREE_MODE = 040000
)

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

// git hash-object filepath.txt
func HashObject(argvs []string, obj_type, _cont string) string {
	if !common.CheckShitInit() {
		fmt.Println("fatal: not a shit repository (or any of the parent directories): .shit")
		return ""
	}

	pwd, _ := os.Getwd()

	var content string
	var filepath string = ""

	if len(_cont) > 0 {
		content = _cont
	} else {
		filepath = argvs[0]

		if argvs[0] == "-w" {
			filepath = pwd + "/" + argvs[1]
		} else if strings.Split(argvs[0], "")[0] == "./" {
			filepath = argvs[0]
		} else {
			filepath = pwd + "/" + argvs[0]
		}

		content = common.ReadFileContents(filepath)
	}
	header := fmt.Sprintf("%s %d::\u0000", obj_type, len(content))

	store := header + content
	h := sha1.New()
	io.WriteString(h, store)

	sha1 := fmt.Sprintf("%x", h.Sum(nil))

	dir_p := pwd + "/.shit/objects/" + sha1[0:2]
	f_path := sha1[2:]

	var comp_content bytes.Buffer

	zlibwriter := zlib.NewWriter(&comp_content)
	zlibwriter.Write([]byte(store))
	zlibwriter.Close()

	if !common.CheckDirExistence(dir_p) {
		common.HandleCreateDir(dir_p, 0755)
	}

	if !common.CheckDirExistence(dir_p + "/" + f_path) {
		common.HandleCreateFile(dir_p, f_path)
	}

	err := os.WriteFile(dir_p+"/"+f_path, (comp_content.Bytes()), 0644)

	if err != nil {
		panic("Cannot write hash-object")
	}

	fmt.Println(sha1)
	return sha1
}

func WriteTree() string {
	sha_1 := HashObject([]string{"/.shit/index"}, "tree", "")
	fmt.Println(sha_1)

	return sha_1
}

func CatFile(argvs []string) (string, error) {
	filepath := fmt.Sprintf("./.shit/objects/%s/%s", argvs[1][0:2], argvs[1][2:])

	fbytes, _ := os.ReadFile(filepath)
	buff := bytes.NewBuffer(fbytes)
	r, err := zlib.NewReader(buff)

	if err != nil {
		return "", err
	}

	io.Copy(buff, r)
	store := buff.String()

	return store, nil
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
		sha_1 = HashObject([]string{filename}, "blob", "")
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

// [<tree_id>, |< -p > <sha_1>|, < -m > "commit message"]
// If the index hasn't changed, that is, an sha_1 already exists, don't commit, show approp. message instead
func CommitTree(treesha_1, parentsha_1, commit_msg string) string {
	var commit_content string
	if len(parentsha_1) < 1 {
		commit_content = fmt.Sprintf("tree %s\nauthor %s <%s> %d\n\t%s", treesha_1, "shoebilyas123", "shoebilyas123@gmail.com", time.Now().UnixMilli(), commit_msg)
	} else {
		commit_content = fmt.Sprintf("tree %s\nparent %s\nauthor %s <%s> %d\n\t%s", treesha_1, parentsha_1, "shoebilyas123", "shoebilyas123@gmail.com", time.Now().UnixMilli(), commit_msg)
	}

	commit_hash := HashObject([]string{}, "tree", commit_content)
	fmt.Println(commit_hash)
	return commit_hash
}

func UpdateRef(path string, commit_sha string) bool {
	path_arr := strings.Split(path, "/")
	dir_path := "./.shit/" + strings.Join(path_arr[0:len(path_arr)-1], "/")
	filename := path_arr[len(path_arr)-1]

	fmt.Printf("%s::%s\n", dir_path, filename)

	if !common.CheckDirExistence(path) {
		common.HandleCreateFile(dir_path, filename)
	}

	os.WriteFile("./.shit/"+path, []byte(commit_sha), 0644)

	return true

}

// shit log <ref_name>
func Log(ref, commit_sha string) {
	if len(ref) > 0 {
		if !common.CheckDirExistence("./.shit/" + ref) {
			fmt.Println("Fatal: ref does not exist")
			return
		}

		commit_id, _ := os.ReadFile("./.shit/" + ref)
		Log("", string(commit_id))
	} else if len(commit_sha) > 0 {
		store, _ := CatFile([]string{"", commit_sha})
		parent := (strings.Split(strings.Split(store, "::")[1], "\n")[1])

		fmt.Printf("Commit: %s\n%s\n\n-----\n", commit_sha, strings.Split(store, "::")[1])
		if strings.Split(parent, " ")[0] == "parent" {
			Log("", strings.Split(parent, " ")[1])
		}
	}

}

func CheckoutBranch(branch_name string) {

}
