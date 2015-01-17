/*
Safe rm is a command line utility which tries to make the "rm" command safer.

Running "saferm <path>" will not delete a file, but will instead move it into a ~/.safetrash/
directory.

The .safetrash contains a configuration file called .trashconfig, which lists the files in the
.safetrash. The .trashconfig also contains the TRASH SIZE of the .safetrash (this is the first
number in the .trashconfig). When the contents of the .safetrash together exceed this size, the
oldest contents of the .safetrash will be deleted until the actual size of the .safetrash falls
below the specified trash size (this is not implemented yet).

Currently, the usage of saferm is as follows:

saferm <path>
Move the FILE at the path to the .safetrash.

saferm -trashsize <MB>
Change the trash size to the given number of megabytes.

Next, I will implement

saferm -r <path>
If <path> is a file, it will be moved to the .safetrash. If it is a directory, it will be
recursively moved to the .safetrash.

saferm -cleartrash
Permanently delete the contents of the .safetrash.

saferm -contents
List the contents of the .safetrash.

*/
package main

import (
	"flag"                                       // For command line args.
	"github.com/hariharsubramanyam/saferm/trash" // For trash operations.
)

func main() {
	trashSize := flag.Int64("trashsize", -1, "Set the trash size in MB")
	flag.Parse()

	if *trashSize != -1 { // Attempt to set the trash size.
		SetTrashSize(trashSize)
	} else if flag.NArg() > 0 { // Attempt to delete the file at the path.
		Delete(flag.Arg(0))
	}

}

func Delete(path string) {
	userTrash := trash.NewTrash()
	userTrash.DeleteFile(path)
	userTrash.Save()
}

func SetTrashSize(trashSize *int64) {
	userTrash := trash.NewTrash()
	if *trashSize >= trash.MinTrashSize && *trashSize <= trash.MaxTrashSize {
		userTrash.TrashSize = *trashSize
		userTrash.Save()
	}
}
