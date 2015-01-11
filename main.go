/*
Safe rm

Run as saferm <filename>

Moves file to .safetrash in home directory.
*/
package main

import (
	"github.com/hariharsubramanyam/saferm/trash"
	"os"
)

const TRASHNAME = ".safetrash"

func main() {
	var userTrash *trash.Trash = trash.NewTrash()
	workingDirectory, err := os.Getwd()
	if err == nil {
		fileName := os.Args[1]
		userTrash.DeleteFile(workingDirectory, fileName)
	}
	userTrash.Save()
}
