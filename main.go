/*
Safe rm

Run as saferm <filename>

This will move the given file into the .safetrash.
*/
package main

//TODO(harihar): Update the package comment after the MVP.

import (
	"flag"
	"fmt"
	"github.com/hariharsubramanyam/saferm/trash"
	"os"
)

func main() {
	trashSize := flag.Int("trashsize", -1, "Set the trash size in MB")

	flag.Parse()
	// Create a Trash object to handle the user's action.
	userTrash := trash.NewTrash()

	newTrashSize := 1024 * 1024 * *trashSize
	if newTrashSize >= trash.MinTrashSize && newTrashSize <= trash.MaxTrashSize {
		fmt.Println("Changing trash size")
		userTrash.TrashSize = newTrashSize
		userTrash.Save()
	}

	workingDirectory, err := os.Getwd()
	if err == nil {
		// Delet the file.
		fileName := flag.Arg(0)
		fmt.Println(fileName)
		userTrash.DeleteFile(workingDirectory, fileName)
	}

	// Update the config file.
	userTrash.Save()
}
