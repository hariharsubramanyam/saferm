/*
Safe rm

Run as saferm <filename>

Moves file to .safetrash in home directory.
*/
package main

import (
	"os"
	"path"
	"runtime"
)

const TRASHNAME = ".safetrash"

func main() {

	workingDirectory, err := os.Getwd()
	if err == nil {
		fileName := os.Args[1]
		moveFileToTrash(fileName, workingDirectory)
	}
}

func moveFileToTrash(fileName string, containingDirectory string) {
	originalPath := path.Join(containingDirectory, fileName)
	if pathExists(originalPath) {
		// Create ~/.safetrash if it doesn't exist.
		trashPath := path.Join(homeDirectoryPath(), TRASHNAME)
		if !pathExists(trashPath) {
			os.Mkdir(trashPath, os.ModePerm)
		}

		newPath := path.Join(trashPath, fileName)
		os.Rename(originalPath, newPath)
	}
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func homeDirectoryPath() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
