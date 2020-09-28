package util

import (
	"os"
	"path/filepath"
)

// EnsureDir checks a file could be written to a path, creates the directories as needed
func EnsureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
}

// Join joins together a path and filename
func Join(path, filename string) string {
	return path + string(os.PathSeparator) + filename
}

// FileExists reports whether the named file or directory exists.
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// MinUint32 is a helper function to return the minimum of two uint32s. This avoids a math import and the need to cast
// to floats.
func MinUint32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}
