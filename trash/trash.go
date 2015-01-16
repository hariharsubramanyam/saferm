/*
Package trash contains logic for moving files to the .safetrash, deleting them permanently, and
updating the config file. Most of this funcionality is in the Trash struct.
*/
package trash

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

// Trash is the object representing useful info about the .safetrash.
type Trash struct {
	TrashSize    int      // The size of the .safetrash in bytes.
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
	t.TrashPath = path.Join(HomeDirectoryPath(), TrashDirectoryName)
	t.ConfigPath = path.Join(t.TrashPath, ConfigFileName)
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
			trashSize, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return t
			}
			t.TrashSize = trashSize
			if t.TrashSize > MaxTrashSize {
				t.TrashSize = MaxTrashSize
			} else if t.TrashSize < MinTrashSize {
				t.TrashSize = MinTrashSize
			}
		}

		for scanner.Scan() {
			deletedItem := scanner.Text()
			t.DeletedItems = append(t.DeletedItems, deletedItem)
		}
	}
	return t
}

// DeleteFile moves a file (fileName) inside a directory (containingDir) into the .safetrash.
func (t *Trash) DeleteFile(containingDir string, fileName string) {
	originalPath := path.Join(containingDir, fileName)
	if PathExists(originalPath) {
		newPath := path.Join(t.TrashPath, fileName)
		os.Rename(originalPath, newPath)
		t.DeletedItems = append(t.DeletedItems, fileName)
	}
}

// Save updates the .trashconfig, with the current values stored in the Trash object.
func (t *Trash) Save() {
	if !PathExists(t.TrashPath) {
		os.Mkdir(t.TrashPath, os.ModePerm)
	}
	configString := strconv.Itoa(t.TrashSize)
	for _, deletedItem := range t.DeletedItems {
		configString += "\n" + deletedItem
	}
	ioutil.WriteFile(t.ConfigPath, []byte(configString), os.ModePerm)
}
