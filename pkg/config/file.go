package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"rope/pkg/helpers"
)

//LoadFile load file from current working directory
func LoadFile() ([]byte, error) {
	file, err := ioutil.ReadFile(path.Join(helpers.ProjectDir, ".rope.yaml"))
	if err != nil {
		return nil, err
	}
	return file, err
}

//File config file struct
type File struct {
	Services map[string]int `yaml:"services"`
}

//ParseConfig decode input into config file struct
func ParseConfig(input []byte) (File, error) {
	var output File
	if err := yaml.Unmarshal(input, &output); err != nil {
		return File{}, err
	}
	return output, nil
}
