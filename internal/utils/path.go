package utils

import (
	"os"
)

// IsDir tells a path is a directory or not
func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.Mode().IsDir()
}
