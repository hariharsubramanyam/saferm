#Safe `rm`

Safe `rm`is a command line utility written in Golang which moves files to a `.safetrash/` directory instead of deleting them. When the `.safetrash/` directory gets too full, the oldest item in `.safetrash/` will be deleted.

#Usage

`saferm [options] <path>`

The options are:

`-setsize` - Set the size of the `.safetrash` in MB.

`-contents` - List the contents of the `.safetrash`.

`-cleartrash` - Delete everything in the `.safetrash`.

`-verbose` - Print details of what `saferm` is doing.

`-used` - Print the space used and the total size of the `safetrash` (ex. `30 MB out of 50 MB`).

`-r` - Recursive deletes (for directories).
