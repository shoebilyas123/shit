package object

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/shoebilyas123/shit/common"
)

// git hash-object filepath.txt
func HashObject(argvs []string, obj_type string) string {
	if !common.CheckShitInit() {
		fmt.Println("fatal: not a shit repository (or any of the parent directories): .shit")
		return ""
	}

	pwd, _ := os.Getwd()

	filepath := argvs[0]

	if argvs[0] == "-w" {
		filepath = pwd + "/" + argvs[1]
	} else if strings.Split(argvs[0], "")[0] == "./" {
		filepath = argvs[0]
	} else {
		filepath = pwd + "/" + argvs[0]
	}

	content := common.ReadFileContents(filepath)
	header := fmt.Sprintf("%s %d::\u0000", obj_type, len(content))

	store := header + content
	h := sha1.New()
	io.WriteString(h, store)

	sha1 := fmt.Sprintf("%x", h.Sum(nil))

	dir_p := pwd + "/.shit/objects/" + sha1[0:2]
	f_path := sha1[2:]

	var comp_content bytes.Buffer

	zlibwriter := zlib.NewWriter(&comp_content)
	zlibwriter.Write([]byte(sha1))
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
	sha_1 := HashObject([]string{"/.shit/index"}, "tree")
	fmt.Println(sha_1)

	return sha_1
}
