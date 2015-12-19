package build

import (
	"fmt"
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
	if name == "" {
		cf, err := GetFile(defaultFileNameYaml)
		if err == nil {
			return cf, nil
		}
		cf, err = GetFile(defaultFileNameYml)
		if err == nil {
			return cf, nil
		}
		return nil, fmt.Errorf("neither %s nor %s exists", defaultFileNameYaml, defaultFileNameYml)
	}
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	cf := &File{}
	if err := yaml.Unmarshal(b, cf); err != nil {
		return nil, err
	}
	return cf, nil
}

// GetFileOrDie calls GetFile and if it returned an error, logs and exits
func GetFileOrDie(name string) *File {
	cf, err := GetFile(name)
	if err != nil {
		log.Die("build file not found (%s)", err)
		return nil
	}
	return cf
}

type File struct {
	Version      int           `yaml:"version"`
	Vars         []Var         `yaml:"vars"`
	StepIncludes []StepInclude `yaml:"steps"`
	Pipelines    []Pipeline    `yaml:"pipelines"`
}
