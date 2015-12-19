package build

import (
	"io/ioutil"

	"github.com/arschles/gci/log"
	"gopkg.in/yaml.v2"
)

const (
	defaultFileNameYaml = "gci.yaml"
	defaultFileNameYml  = "gci.yml"
)

// GetFile attempts to get and decode the File at name. If name is empty,
// tries defaultFileNameYaml and then defaultFileNameYml. If no file at name exists,
// or name was empty and neither defaultFileNameYaml nor defaultFileNameYml exists,
// returns ErrNoFile
func GetFile(name string) (*File, error) {
	var fileBytes []byte
	var err error
	fileNames := []string{name, defaultFileNameYaml, defaultFileNameYml}
	for _, fileName := range fileNames {
		b, err := ioutil.ReadFile(fileName)
		if err == nil {
			fileBytes = b
			break
		}
	}
	if err != nil {
		return nil, err
	}
	cf := &File{}
	if err := yaml.Unmarshal(fileBytes, cf); err != nil {
		return nil, err
	}
	return cf, nil
}

// GetFileOrDie calls GetFile and if it returned an error, logs and exits
func GetFileOrDie(name string) *File {
	cf, err := GetFile(name)
	if err != nil {
		log.Die(err.Error())
		return nil
	}
	return cf
}

type File struct {
	Version      string        `yaml:"version"`
	Vars         []Var         `yaml:"vars"`
	StepIncludes []StepInclude `yaml:"steps"`
	Pipelines    []Pipeline    `yaml:"pipelines"`
}
