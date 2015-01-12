/*
Safe rm

Run as saferm <filename>

This will move the given file into the .safetrash.
*/
package main

//TODO(harihar): Update the package comment after the MVP.

import (
	"github.com/hariharsubramanyam/saferm/trash"
	"os"
)

func main() {
	// Create a Trash object to handle the user's action.
	userTrash := trash.NewTrash()

	workingDirectory, err := os.Getwd()
	if err == nil {
		// Delet the file.
		fileName := os.Args[1]
		userTrash.DeleteFile(workingDirectory, fileName)
	}

	// Update the config file.
	userTrash.Save()
}
