package utils

import (
	"os"
	"path/filepath"
)

func tempUploadDir() string {
	dir := filepath.Join("upload", "temp")

	// auto create kalau belum ada
	_ = os.MkdirAll(dir, 0755)

	return dir
}

func TempFilePath(filename string) string {
	return filepath.Join(tempUploadDir(), filename)
}
