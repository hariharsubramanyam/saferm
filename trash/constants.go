package trash

const (
	MaxTrashSize       = 1024           // The trash cannot exceed 10 GB in size.
	MinTrashSize       = 1              // The trash size cannot be less than 1 MB.
	DefaultTrashSize   = 10             // By default, trash is 10 MB large.
	TrashDirectoryName = ".safetrash"   // The trash directory is called .safetrash.
	ConfigFileName     = ".trashconfig" // The config file is called .trashconfig.
)
