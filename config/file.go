package config

import (
	"io/ioutil"

	"github.com/arschles/gci/config/ci"
	"github.com/arschles/gci/log"
	"gopkg.in/yaml.v2"
)

const (
	DefaultFileNameYaml = "gci.yaml"
	DefaultFileNameYml  = "gci.yml"
)

func ReadBytes(b []byte) (*File, error) {
	var cf File
	if err := yaml.Unmarshal(b, &cf); err != nil {
		return nil, err
	}
	return &cf, nil
}

// GetFile attempts to get and decode the File at name. If name is empty,
// tries DefaultFileNameYaml and then defaultFileNameYml. If no file at name exists,
// or name was empty and neither DefaultFileNameYaml nor defaultFileNameYml exists,
// returns ErrNoFile
func Read(name string) (*File, error) {
	var fileBytes []byte
	var err error
	fileNames := []string{name, DefaultFileNameYaml, DefaultFileNameYml}
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
	return ReadBytes(fileBytes)
}

// ReadOrDie calls Read and if it returned an error, logs and exits
func ReadOrDie(name string) *File {
	cf, err := Read(name)
	if err != nil {
		log.Die("Reading config file %s [%s]", name, err)
		return nil
	}
	return cf
}

func Empty() *File {
	return &File{}
}

type File struct {
	Version string         `yaml:"version"`
	Build   Build          `yaml:"build"`
	Test    Test           `yaml:"test"`
	Docker  Docker         `yaml:"docker"`
	Custom  []CustomTarget `yaml:"custom"`
	CI      ci.CI          `yaml:"ci"`
}

func (f File) String() string {
	return "GCI Config file version " + f.Version
}

type Build struct {
	OutputBinary string   `yaml:"output-binary"`
	CrossCompile bool     `yaml:"cross-compile"`
	Env          []string `yaml:"env"`
}

func (b Build) GetOutputBinary(pathBase string) string {
	if b.OutputBinary == "" {
		return pathBase
	}
	return b.OutputBinary
}

type Test struct {
	Paths []string `yaml:"paths"`
	Env   []string `yaml:"env"`
}

func (t Test) GetPaths() []string {
	if len(t.Paths) == 0 {
		return []string{"./..."}
	}
	return t.Paths
}
