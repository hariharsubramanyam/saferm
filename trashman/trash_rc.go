package trashman

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

const configfile = ".trashconfig"

type Config struct {
	TrashSize int
}

func (c *Config) WriteConfig() {
	trashPath := TrashPath()
	trashSize := strconv.Itoa(c.TrashSize)
	configPath := path.Join(trashPath, configfile)
	ioutil.WriteFile(configPath, []byte(trashSize), os.ModePerm)
}

func NewConfig() (*Config, error) {
	c := &Config{TrashSize: DefaultTrashSize}
	configPath := path.Join(TrashPath(), configfile)
	if PathExists(configPath) {
		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			return c, err
		}
		trashSize, err := strconv.Atoi(string(data))
		if err != nil {
			return c, err
		}
		c.TrashSize = trashSize
	}
	return c, nil
}
