package main

import (
	"io/ioutil"

	"github.com/arschles/gocons/log"
	"gopkg.in/yaml.v2"
)

const defaultConsFileName = "gocons.yaml"

func getConsfile() (*Consfile, error) {
	b, err := ioutil.ReadFile(defaultConsFileName)
	if err != nil {
		return nil, err
	}
	f := &Consfile{}
	if err := yaml.Unmarshal(b, f); err != nil {
		return nil, err
	}
	return f, nil
}

func getConsfileOrDie() *Consfile {
	consfile, err := getConsfile()
	if err != nil {
		log.Die("error getting consfile [%s]", err)
		return nil
	}
	return consfile
}

type Consfile struct {
	Version int      `yaml:"version"`
	Plugins []string `yaml:"repos"`
	Targets []Target `yaml:"targets"`
}

type Target struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Depends     string   `yaml:"depends"`
	Commands    []string `yaml:"commands"`
}
