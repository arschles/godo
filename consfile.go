package main

import (
	"fmt"
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
	Envs    []Env    `yaml:"environment_vars"`
	Plugins []string `yaml:"repos"`
	Targets []Target `yaml:"targets"`
}

type Env struct {
	Name string `yaml:"name"`
	Val  string `yaml:"val"`
}

func (e Env) String() string {
	return fmt.Sprintf("%s=%s", e.Name, e.Val)
}

type Envs []Env

func (e Envs) Strings() []string {
	strs := make([]string, len(e))
	for i, env := range e {
		strs[i] = env.String()
	}
	return strs
}

type Target struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Depends     string   `yaml:"depends"`
	Commands    []string `yaml:"commands"`
}
