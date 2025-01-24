package fileutil

import (
	"os"
	"path/filepath"
	"strings"
)

func ConvertPath(path string, from string, to string) string {
	return strings.Replace(path, from, to, 1)
}

func EnsureDirectory(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0755)
}

func ChangeExtension(path string, newExt string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext) + newExt
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
