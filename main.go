/*

Running "saferm <path>" will not delete a file, but will instead move it into a ~/.safetrash/
directory.

Safe `rm`is a command line utility written in Golang which moves files to a `.safetrash/`
directory instead of deleting them. When the `.safetrash/` directory gets too full,
the oldest item in `.safetrash/` will be deleted.

The .safetrash contains a configuration file called .trashconfig, which lists the files in the
.safetrash. The .trashconfig also contains the TRASH SIZE of the .safetrash (this is the first
number in the .trashconfig). When the contents of the .safetrash together exceed this size, the
oldest contents of the .safetrash will be deleted until the actual size of the .safetrash falls
below the specified trash size (this is not implemented yet).

`saferm [options] <path>`

The options are:

`-setsize` - Set the size of the `.safetrash` in MB.

`-contents` - List the contents of the `.safetrash`.

`-cleartrash` - Delete everything in the `.safetrash`.

`-verbose` - Print details of what `saferm` is doing.

`-used` - Print the space used and the total size of the `safetrash` (ex. `30 MB out of 50 MB`).

`-r` - Recursive deletes (for directories).
*/
package main

import (
	"flag"                                       // For command line args.
	"fmt"                                        // For printing trash contents.
	"github.com/hariharsubramanyam/saferm/trash" // For trash operations.
	"strings"                                    // For joining contents slice into string.
)

func main() {
	trashSize := flag.Int64("setsize", -1, "Set the trash size in MB")
	contents := flag.Bool("contents", false, "Display the contents of the .safetrash")
	clearTrash := flag.Bool("cleartrash", false, "Delete everything in the .safetrash")
	verbose := flag.Bool("verbose", false, "Print verbose output during trash operations")
	used := flag.Bool("used", false, "See the space used in the trash and its total size in MB")
	recursive := flag.Bool("r", false, "Recursive delete (for directories)")
	flag.Parse()

	userTrash := trash.NewTrash()

	if *verbose {
		// Make verbose outputs.
		userTrash.Verbose = true
	}

	if *contents {
		// Print the contents of the trash.
		fmt.Println("Trash Contents", strings.Join(userTrash.Contents(), ", "))
	}

	if *trashSize >= trash.MinTrashSize && *trashSize <= trash.MaxTrashSize {
		// Reset the trash size.
		userTrash.TrashSize = *trashSize
	}

	if *clearTrash {
		// Delete all the items in the trash.
		userTrash.ClearTrash()
	}

	if *used {
		fmt.Println("Used Space: ", userTrash.SpaceUsed()/1024/1024, "MB, Capacity:",
			userTrash.TrashSize, "MB")
	}

	if flag.NArg() > 0 {
		path := flag.Arg(0)
		if trash.PathExists(path) {
			if *recursive {
				fmt.Println("Here i am")
				userTrash.Delete(path)
			} else {
				userTrash.DeleteFile(path)
			}
		}
	}

	userTrash.Save()
}
