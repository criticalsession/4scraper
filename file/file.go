package file

import (
	"os"
	"path/filepath"
)

func FileExists(dir, file string) bool {
	p := dir
	if file != "" {
		p = filepath.Join(dir, file)
	}

	_, err := os.Stat(p)
	return !os.IsNotExist(err)
}

func DirExists(dir string) bool {
	return FileExists(dir, "")
}

func GetExtension(filename string) string {
	return filepath.Ext(filename)
}

func SplitFilename(filename string) (baseName string, ext string) {
	ext = GetExtension(filename)
	baseName = filename[:len(filename)-len(ext)]

	return
}
