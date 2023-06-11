package files

import (
	"embed"
	"fmt"
)

func GetAll(fs embed.FS, fileName string) []byte {
	jsonFileBytes, err := fs.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}
	return jsonFileBytes
}
