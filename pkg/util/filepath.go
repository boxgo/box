package util

import (
	"path/filepath"
	"strings"
)

func Filename(filePath string) string {
	_, file := filepath.Split(filePath)
	ext := filepath.Ext(filePath)
	return strings.Replace(file, ext, "", 1)
}
