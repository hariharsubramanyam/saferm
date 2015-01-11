package trash

import (
	"os"      // For file ops.
	"runtime" // For getting OS version (to find HOME directory).
)

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

// PathExists determines whether a given path (to either a file or directory) exists.
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
