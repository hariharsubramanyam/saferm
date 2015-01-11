/*
Package trashman implements code for managing files for the moving (and eventually deleting) files
in the safe trash.
*/
package trashman

import (
	"os"      // For file ops.
	"path"    // For joining paths.
	"runtime" // For getting OS version (to find HOME directory).
)

// trashfile is the directory where "deleted" files/directories will be moved to.
const trashfile = ".safetrash"

// HomeDirectoryPath is the path to the user's HOME directory (or equivalent, if not UNIX).
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

// TrashPath is the path to the trash where "deleted" files go to.
func TrashPath() string {
	return path.Join(HomeDirectoryPath(), trashfile)
}

// PathExists determines whether a given path (to either a file or directory) exists.
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// MoveFileToTrash moves the given fileName (which is inside the given containingDirectory) to the
// safe trash.
func MoveFileToTrash(containingDirectory string, fileName string) {
	originalPath := path.Join(containingDirectory, fileName)
	if !PathExists(originalPath) {
		os.Mkdir(TrashPath(), os.ModePerm)
	}
	newPath := path.Join(TrashPath(), fileName)
	os.Rename(originalPath, newPath)
}
