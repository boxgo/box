package fputil

import (
	"os"
	"path/filepath"
	"strings"
)

func GetFilename(filePath string) string {
	_, file := filepath.Split(filePath)
	ext := filepath.Ext(filePath)
	return strings.Replace(file, ext, "", 1)
}

func FirstExistFilePath(paths []string) string {
	for _, p := range paths {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			return p
		}
	}

	return ""
}

// IsFile checks whether the path is a file,
// it returns false when it's a directory or does not exist.
func IsFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}

	return !fi.IsDir()
}
