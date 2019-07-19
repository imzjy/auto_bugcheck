package main

import (
	"os"
)

// FileExist check file or folder exist, return true if exist
func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
