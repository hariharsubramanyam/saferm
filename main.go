/*
Safe rm

Run as saferm <filename>

Moves file to .safetrash in home directory.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/hariharsubramanyam/saferm/trashman"
	"os"
)

const TRASHNAME = ".safetrash"

func main() {
	recursive := flag.Bool("r", false, "Whether the delete should be recursive")
	flag.Parse()
	fmt.Println(*recursive)

	workingDirectory, err := os.Getwd()
	if err == nil {
		fileName := flag.Arg(0)
		fmt.Println("Deleting", fileName)
		trashman.MoveFileToTrash(workingDirectory, fileName)
	}
}
