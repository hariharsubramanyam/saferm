/*
Package trash contains logic for moving files to the .safetrash, deleting them permanently, and
updating the config file. Most of this funcionality is in the Trash struct.
*/
package trash

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// Trash is the object representing useful info about the .safetrash.
type Trash struct {
	TrashSize    int64    // The size of the .safetrash in bytes.
	TrashPath    string   // The path of the .safetrash (inside the HOME directory).
	ConfigPath   string   // The path of the .trashconfig file (inside .safetrash directory).
	DeletedItems []string // The slice of items that have been deleted (most recent is last).
}

// NewTrash creates a Trash object (reading from the .trashconfig, if it exists).
// If there is no .trashconfig, then the .safetrash begins with a default capacity
// (see the constants.go file).
func NewTrash() *Trash {
	t := &Trash{}

	// Set the paths and default size.
	t.TrashPath = filepath.Join(HomeDirectoryPath(), TrashDirectoryName)
	t.ConfigPath = filepath.Join(t.TrashPath, ConfigFileName)
	t.TrashSize = DefaultTrashSize

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

// DeleteFile deletes the file at the given path. If the path points to a directory, this function
// does nothing.
func (t *Trash) DeleteFile(path string) {
	// First get an absolute path.
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}

	// Ensure that this is a valid path and NOT a directory.
	if isDir, err := IsDirectory(absPath); PathExists(absPath) && err == nil && !isDir {
		// Move the file from its current location to .safetrash/
		fileName := filepath.Base(absPath)
		newPath := filepath.Join(t.TrashPath, fileName)
		os.Rename(absPath, newPath)
		t.DeletedItems = append(t.DeletedItems, fileName)
		t.DeleteOldestIfNeeded()
	}
}

// DeleteOldestIfNeeded keeps deleting the oldest item from the trash until its size is
// smaller than the permitted size.
func (t *Trash) DeleteOldestIfNeeded() {
	usedSpace := t.SpaceUsed()
	trashSizeInBytes := t.TrashSize * 1024 * 1024
	lastDeletedIndex := -1
	for usedSpace > trashSizeInBytes {
		for i, deletedItem := range t.DeletedItems {
			pathToDeletedItem := filepath.Join(t.TrashPath, deletedItem)
			if PathExists(pathToDeletedItem) {
				os.Remove(pathToDeletedItem)
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

func (t *Trash) Contents() []string {
	contents := make([]string, 0)
	updateContents := func(path string, info os.FileInfo, err error) error {
		base := filepath.Base(path)
		if base != ConfigFileName && base != TrashDirectoryName {
			contents = append(contents, base)
		}
		return nil
	}
	filepath.Walk(t.TrashPath, updateContents)
	return contents
}

func (t *Trash) ClearTrash() {
	contents := t.Contents()
	for _, content := range contents {
		os.Remove(filepath.Join(t.TrashPath, content))
	}
	t.DeletedItems = make([]string, 0)
}

// Save updates the .trashconfig, with the current values stored in the Trash object.
func (t *Trash) Save() {
	// Create the .safetrash/ if it doesn't exist.
	if !PathExists(t.TrashPath) {
		os.Mkdir(t.TrashPath, os.ModePerm)
	}

	// Write the .trashconfig file.
	configString := strconv.FormatInt(t.TrashSize, 10)
	for _, deletedItem := range t.DeletedItems {
		configString += "\n" + deletedItem
	}
	ioutil.WriteFile(t.ConfigPath, []byte(configString), os.ModePerm)
}
