package config

import (
	"io/ioutil"

	"github.com/arschles/godo/log"
	"gopkg.in/yaml.v2"
)

const (
	DefaultFileNameYaml = "godo.yaml"
	DefaultFileNameYml  = "godo.yml"
)

func ReadBytes(b []byte) (*File, error) {
	var cf File
	if err := yaml.Unmarshal(b, &cf); err != nil {
		return nil, err
	}
	return &cf, nil
}

// Read attempts to get and decode the File at name. If name is empty,
// tries DefaultFileNameYaml and then defaultFileNameYml. If no file at name exists,
// or name was empty and neither DefaultFileNameYaml nor defaultFileNameYml exists,
// returns ErrNoFile
func Read(name string) (*File, error) {
	var fileBytes []byte
	var err error
	fileNames := []string{name, DefaultFileNameYaml, DefaultFileNameYml}
	for _, fileName := range fileNames {
		b, readFileErr := ioutil.ReadFile(fileName)
		if readFileErr == nil {
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

// Empty returns an empty config file
func Empty() *File {
	return &File{}
}

// File is the complete in-memory representation of a config file
type File struct {
	Version string         `yaml:"version"`
	Build   Build          `yaml:"build"`
	Docker  Docker         `yaml:"docker"`
	Custom  []CustomTarget `yaml:"custom"`
}

func (f File) String() string {
	return "Godo Config file version " + f.Version
}

// Build is the configuration for a build
type Build struct {
	ImageName    string   `yaml:"image-name"`
	ImageTag     string   `yaml:"image-tag"`
	Gopath       string   `json:"gopath"`
	OutputBinary string   `yaml:"output-binary"`
	Env          []string `yaml:"env"`
}

func (b Build) GetOutputBinary(pathBase string) string {
	if b.OutputBinary == "" {
		return pathBase
	}
	return b.OutputBinary
}
