package trash

import (
	"os" // For file ops.
)

// HomeDirectoryPath is the path to the user's HOME directory (or equivalent, if not UNIX).
func HomeDirectoryPath() string {
	return os.Getenv("HOME")
}

// PathExists determines whether a given path (to either a file or directory) exists.
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
