package trash

import (
	"os" // For file ops.
)

func HomeDirectoryPath() string {
	return os.Getenv("HOME")
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}
