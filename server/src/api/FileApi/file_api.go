package FileApi

import (
	"os"
)

type FileApi struct {
}

func ChkExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
