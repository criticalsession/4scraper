package main

import (
	"os"
	"path/filepath"
)

func FileExists(dir, filename string) bool {
	path := dir
	if filename != "" {
		path = filepath.Join(dir, filename)
	}

	_, err := os.Stat(path)
	if os.IsNotExist(err) { // file doesn't exist
		return false
	} else {
		return true
	}
}

func DirExists(dir string) bool {
	return FileExists(dir, "")
}

func GetExtension(filename string) string {
	return filepath.Ext(filename)
}
