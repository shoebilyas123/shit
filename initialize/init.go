package initialize

import (
	"fmt"
	"io/fs"
	"os"
)

func CheckDirExistence(path string) bool {
	if _, err := os.Open(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func HandleCreateDir(path string, perm fs.FileMode) {
	err := os.Mkdir(path, 0755)

	if os.IsExist(err) {
		panic("Error: Cannot create the target directory")
	}
}

func Init(argvs []string) {
	pwd, err := os.Getwd()

	if err != nil {
		panic("Error: Cannot get the present working directory.")
	}
	target_dir := pwd

	// Check if the argvs is empty or .
	if len(argvs) > 0 && argvs[0] != "." {
		target_dir = pwd + "/" + argvs[0]
	}

	// If a name is provided create that named directory

	if len(argvs) > 0 && argvs[0] != "." {
		if !CheckDirExistence(target_dir) {
			HandleCreateDir(target_dir, 0755)
		} else {
			// TODO: Will ask if you want to init that directory?
			// Check for directory name clashes
			fmt.Printf("Navigate into the directory and shit init your project")
		}
	}

	// If the target_dir already has a .shit folder throw appropriate error
	if CheckDirExistence(target_dir + "/.shit") {
		panic("Error: Cannot overwrite an already initialized <Repository>")
	}

	// In that directory create a .shit directory
	HandleCreateDir(target_dir+"/.shit", 0755)
	// In that .shit directory create an objs directory
	HandleCreateDir(target_dir+"/.shit/objs", 0755)
}
