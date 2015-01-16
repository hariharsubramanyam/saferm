/*
Safe rm

Run as saferm <filename>

This will move the given file into the .safetrash.
*/
package main

//TODO(harihar): Update the package comment after the MVP.

import (
	"flag"
	"github.com/hariharsubramanyam/saferm/trash"
	"os"
)

func main() {
	trashSize := flag.Int("trashsize", -1, "Set the trash size in MB")

	flag.Parse()
	// Create a Trash object to handle the user's action.
	userTrash := trash.NewTrash()

	if *trashSize >= trash.MinTrashSize && *trashSize <= trash.MaxTrashSize {
		userTrash.TrashSize = *trashSize
		userTrash.Save()
	}

	workingDirectory, err := os.Getwd()
	if err == nil {
		// Delet the file.
		fileName := flag.Arg(0)
		userTrash.DeleteFile(workingDirectory, fileName)
	}

	// Update the config file.
	userTrash.Save()
}
