/*
Package trash contains logic for moving files to the .safetrash, deleting them permanently, and
updating the config file. Most of this funcionality is in the Trash struct.
*/
package trash

import (
	"bufio"         // For reading the config file.
	"errors"        // For returning errors if the path is invalid.
	"fmt"           // For verbose logging.
	"io/ioutil"     // For writing the config file.
	"os"            // For delete, rename, directory access, and much more.
	"path/filepath" // For combining and manipulating path strings.
	"strconv"       // For parsing and integers and converting integers to strings.
)

// Trash is the object representing useful info about the .safetrash.
type Trash struct {
	TrashSize    int64    // The size of the .safetrash in bytes.
	TrashPath    string   // The path of the .safetrash (inside the HOME directory).
	ConfigPath   string   // The path of the .trashconfig file (inside .safetrash directory).
	DeletedItems []string // The slice of items that have been deleted (most recent is last).
	Verbose      bool
}

// NewTrash creates a *Trash object (reading from the .trashconfig, if it exists).
// If there is no .trashconfig, then the .safetrash begins with a default capacity
// (see the constants.go file).
func NewTrash() *Trash {
	return NewTrashWithPaths(HomeDirectoryPath(), TrashDirectoryName, ConfigFileName)
}

// NewTrashWithPaths creates a *Trash with its location and names as injected dependencies, this
// can be useful for testing.
func NewTrashWithPaths(dirContainingTrash string, trashName string, configFileName string) *Trash {
	t := &Trash{}

	// Set the paths and default size.
	t.TrashPath = filepath.Join(dirContainingTrash, trashName)
	t.ConfigPath = filepath.Join(t.TrashPath, configFileName)
	t.TrashSize = DefaultTrashSize

	// Create the .safetrash/ if it doesn't exist.
	if !PathExists(t.TrashPath) {
		os.Mkdir(t.TrashPath, os.ModePerm)
	}

	// The trash is not verbose.
	t.Verbose = false

	// Attempt to update size from .trashconfig, if it exists.
	if PathExists(t.ConfigPath) {
		file, err := os.Open(t.ConfigPath)
		if err != nil {
			return t
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			// Read the first line of the .trashconfig, which is the size of the .safetrash.
			trashSize, err := strconv.ParseInt(scanner.Text(), 10, 64)
			if err != nil {
				return t
			}
			t.TrashSize = trashSize

			// Constrain the trash size.
			if t.TrashSize > MaxTrashSize {
				t.TrashSize = MaxTrashSize
			} else if t.TrashSize < MinTrashSize {
				t.TrashSize = MinTrashSize
			}
		}

		// Read the deleted items.
		for scanner.Scan() {
			deletedItem := scanner.Text()
			t.DeletedItems = append(t.DeletedItems, deletedItem)
		}
	}
	return t
}

// Log will print when in verbose mode.
func (t *Trash) Log(messages ...interface{}) {
	if t.Verbose {
		fmt.Println(messages...)
	}
}

// DeleteFile deletes the file at the given path.
// If the path points to a directory, it returns an error.
func (t *Trash) DeleteFile(path string) error {
	// First get an absolute path.
	absPath, err := filepath.Abs(path)
	if err != nil {
		t.Log(err)
		return err
	}

	// Ensure that this is a valid path and NOT a directory.
	if isDir, err := IsDirectory(absPath); err == nil && PathExists(absPath) && !isDir {
		// Move the file from its current location to .safetrash/
		return t.Delete(absPath)
	} else if isDir {
		return errors.New(fmt.Sprintln(absPath, "is a directory, use the -r flag to delete it"))
	} else {
		return errors.New(fmt.Sprintln(absPath, "is not a valid path, are you sure it exists?"))
	}

	return nil
}

// Delete deletes the item (file or directory) at the given absolute path.
func (t *Trash) Delete(absPath string) error {
	fileName := filepath.Base(absPath)
	newPath := filepath.Join(t.TrashPath, fileName)
	os.Rename(absPath, newPath)
	t.DeletedItems = append(t.DeletedItems, fileName)
	t.Log("Deleted", fileName)
	t.DeleteOldestIfNeeded()
	return nil
}

// DeleteOldestIfNeeded keeps deleting the oldest item from the trash until its size is
// smaller than the permitted size.
func (t *Trash) DeleteOldestIfNeeded() {
	usedSpace := t.SpaceUsed()
	trashSizeInBytes := t.TrashSize * 1024 * 1024
	lastDeletedIndex := -1
	for usedSpace > trashSizeInBytes {
		t.Log("Used space:", usedSpace/1024/1024, "MB, Trash size", t.TrashSize, "MB")
		t.Log("Need to clear space...")
		for i, deletedItem := range t.DeletedItems {
			pathToDeletedItem := filepath.Join(t.TrashPath, deletedItem)
			if PathExists(pathToDeletedItem) {
				os.Remove(pathToDeletedItem)
				t.Log("Deleting", deletedItem)
				lastDeletedIndex = i
				break
			} // if
		} // inner for
		usedSpace = t.SpaceUsed()
	} // outer for

	// Update the list of deleted items, because we may have deleted some.
	if lastDeletedIndex != -1 {
		t.DeletedItems = t.DeletedItems[lastDeletedIndex+1:]
	}
} // function

// SpaceUsed determines the current size of the .safetrash.
func (t *Trash) SpaceUsed() int64 {
	var spaceUsed int64 = 0
	updateSpace := func(path string, info os.FileInfo, err error) error {
		spaceUsed += info.Size()
		return nil
	}

	filepath.Walk(t.TrashPath, updateSpace)
	return spaceUsed
}

// Contents returns the files and directories in the .safetrash.
func (t *Trash) Contents() []string {
	paths, err := filepath.Glob(filepath.Join(t.TrashPath, "*"))
	contents := make([]string, 0)
	if err == nil {
		for _, path := range paths {
			base := filepath.Base(path)
			if base == ConfigFileName {
				continue
			}
			contents = append(contents, filepath.Base(base))
		}
	}
	return contents
}

// ClearTrash removes everything in the .safetrash (except for the .trashconfig).
func (t *Trash) ClearTrash() {
	contents := t.Contents()
	for _, content := range contents {
		t.Log("Deleting", content)
		os.Remove(filepath.Join(t.TrashPath, content))
	}
	t.DeletedItems = make([]string, 0)
}

// Save updates the .trashconfig, with the current values stored in the Trash object.
func (t *Trash) Save() {
	// Write the .trashconfig file.
	configString := strconv.FormatInt(t.TrashSize, 10)
	for _, deletedItem := range t.DeletedItems {
		configString += "\n" + deletedItem
	}
	ioutil.WriteFile(t.ConfigPath, []byte(configString), os.ModePerm)
}
