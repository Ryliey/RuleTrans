package fileutil

import (
	"os"
	"path/filepath"
	"strings"
)

func ConvertPath(path string, from string, to string) string {
	normalizedPath := filepath.ToSlash(path)
	from = filepath.ToSlash(from) + "/"
	to = filepath.ToSlash(to) + "/"

	if strings.HasPrefix(normalizedPath, from) {
		return filepath.FromSlash(to + normalizedPath[len(from):])
	}
	return filepath.FromSlash(normalizedPath)
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

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
