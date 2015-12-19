package main

import (
	"fmt"
	"io/ioutil"

	"github.com/arschles/canta/log"
	"gopkg.in/yaml.v2"
)

const (
	defaultConsFileNameYaml = "gocons.yaml"
	defaultConsfileNameYml  = "gocons.yml"
)

func getConsfileBytes() ([]byte, error) {
	b, err := ioutil.ReadFile(defaultConsFileNameYaml)
	if err == nil {
		return b, nil
	}
	b, err = ioutil.ReadFile(defaultConsfileNameYml)
	if err == nil {
		return b, nil
	}
	return nil, err

}

func getConsfile() (*Consfile, error) {
	b, err := getConsfileBytes()
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
	Name    string `yaml:"name"`
	Val     string `yaml:"val"`
	Default string `yaml:"default"`
}

func (e Env) String() string {
	if e.Val == "" {
		return fmt.Sprintf("%s=%s", e.Name, e.Default)
	}
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
