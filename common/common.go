package common

import (
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

func HandleCreateFile(path string, name string) {
	file, err := os.Create(path + "/" + name)

	if err != nil {
		panic("Error: Cannot create file")

	}

	defer file.Close()
}

func ReadFileContents(path string) string {
	contents, err := os.ReadFile(path)

	if err != nil {
		panic("Cannot read file contents")
	}

	return string(contents)
}

func CheckShitInit() bool {
	return CheckDirExistence("./.shit")
}