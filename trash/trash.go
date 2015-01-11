package trash

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type Trash struct {
	TrashSize  int
	TrashPath  string
	ConfigPath string
}

func NewTrash() *Trash {
	t := &Trash{}
	t.TrashPath = path.Join(HomeDirectoryPath(), TrashFileName)
	t.ConfigPath = path.Join(t.TrashPath, ConfigFileName)
	t.TrashSize = DefaultTrashSize
	if PathExists(t.ConfigPath) {
		file, err := os.Open(t.ConfigPath)
		if err != nil {
			return t
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			trashSize, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return t
			}
			t.TrashSize = trashSize
		}
	}
	return t
}

func (t *Trash) DeleteFile(containingDir string, fileName string) {
	originalPath := path.Join(containingDir, fileName)
	if PathExists(originalPath) {
		newPath := path.Join(t.TrashPath, fileName)
		os.Rename(originalPath, newPath)
	}
}

func (t *Trash) Save() {
	if !PathExists(t.TrashPath) {
		os.Mkdir(t.TrashPath, os.ModePerm)
	}
	trashSize := strconv.Itoa(t.TrashSize)
	ioutil.WriteFile(t.ConfigPath, []byte(trashSize), os.ModePerm)
}
