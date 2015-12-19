package config

import (
	"fmt"
	"io/ioutil"

	"github.com/arschles/canta/log"
	"gopkg.in/yaml.v2"
)

const (
	defaultFileNameYaml = "gocons.yaml"
	defaultFileNameYml  = "gocons.yml"
)

// GetFile attempts to get and decode the CantaFile at name. If name is empty,
// tries defaultFileNameYaml and then defaultFileNameYml. If no file at name exists,
// or name was empty and neither defaultFileNameYaml nor defaultFileNameYml exists,
// returns ErrNoFile
func GetFile(name string) (*CantaFile, error) {
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
	cf := &CantaFile{}
	if err := yaml.Unmarshal(b, cf); err != nil {
		return nil, err
	}
	return cf, nil
}

func GetFileOrDie(name string) *CantaFile {
	cf, err := GetFile(name)
	if err != nil {
		log.Die("build file not found (%s)", err)
		return nil
	}
	return cf
}

type CantaFile struct {
	Version int      `yaml:"version"`
	Plugins []Plugin `yaml:"plugins"`
	Vars    []Var    `yaml:"vars"`
	Targets []Target `yaml:"targets"`
}

type Plugin struct {
}

type Var struct {
	Name    string `yaml:"name"`
	Default string `yaml:"default"`
	Env     string `yaml:"env"`
	Value   string `yaml:"value"`
}

type Target struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}
