package utils

import (
	"os"

	"github.com/charmbracelet/log"
)

func GetFileSize(path string) int64 {
	file, err := os.Stat(path)
	if err != nil {
		log.Error("Error getting file stats", err)
		return 0
	}
	return file.Size()
}
