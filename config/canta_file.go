package config

package main

import (
	"fmt"
	"io/ioutil"

	"github.com/arschles/gocons/log"
	"gopkg.in/yaml.v2"
)

const (
	defaultConsFileNameYaml = "gocons.yaml"
	defaultConsfileNameYml  = "gocons.yml"
)

type CantaFile struct {
	Version int      `yaml:"version"`
  Plugins []Plugin `yaml:"plugins"`
  Vars []Var `yaml:"vars"`
  Targets []Target `yaml:"targets"`
}

type Plugin struct {

}

type Var {
  Name string `yaml:"name"`
  Default string `yaml:"default"`
  Env string `yaml:"env"`
  Value string `yaml:"value"`
}

type Target struct {
}
