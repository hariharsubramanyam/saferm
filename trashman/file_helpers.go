/*
Package trashman implements code for managing files for the moving (and eventually deleting) files
in the safe trash.
*/
package trashman

import (
	"os"
	"path"
	"runtime"
)

const trashfile = ".safetrash"

func HomeDirectoryPath() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func TrashPath() string {
	return path.Join(HomeDirectoryPath(), trashfile)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func MoveFileToTrash(containingDirectory string, fileName string) {
	originalPath := path.Join(containingDirectory, fileName)
	if !PathExists(originalPath) {
		os.Mkdir(TrashPath(), os.ModePerm)
	}
	newPath := path.Join(TrashPath(), fileName)
	os.Rename(originalPath, newPath)
}
