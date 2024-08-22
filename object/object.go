package object

import (
	"fmt"

	"github.com/shoebilyas123/shit/common"
)

// git hash-object filepath.txt
func HashObject(filepath string) {
	content := common.ReadFileContents(filepath)
	header := fmt.Sprintf("blob %d\u0000", len(content))

	store := content + header

}
