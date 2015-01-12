#Safe `rm`

Safe `rm`is a command line utility written in Golang which moves files to a `.safetrash/` directory instead of deleting them. When the `.safetrash/` directory gets too full, the oldest item in `.safetrash/` will be deleted.

#Usage

MVP: `saferm <filename>`
Move the file to the .safetrash in the home directory.
**COMPLETE**
