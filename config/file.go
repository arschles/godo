package config

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
func Read(name string) (*File, error) {
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

// ReadOrDie calls Read and if it returned an error, logs and exits
func ReadOrDie(name string) *File {
	cf, err := Read(name)
	if err != nil {
		log.Die("Reading config file %s [%s]", name, err)
		return nil
	}
	return cf
}

type File struct {
	Version     string      `yaml:"version"`
	Build       Build       `yaml:"build"`
	Test        Test        `yaml:"test"`
	DockerBuild DockerBuild `yaml:"docker-build"`
}

type Build struct {
	OutputBinary string `yaml:"output-binary"`
}

func (b Build) GetOutputBinary(pathBase string) string {
	if b.OutputBinary == "" {
		return pathBase
	}
	return b.OutputBinary
}

type Test struct {
	Paths []string `yaml:"paths"`
}

func (t Test) GetPaths() []string {
	if len(t.Paths) == 0 {
		return []string{"./..."}
	}
	return t.Paths
}

type DockerBuild struct {
	ImageName          string `yaml:"image-name"`
	DockerfileLocation string `yaml:"dockerfile-location"`
}

func (d DockerBuild) GetDockerfileLocation() string {
	if d.DockerfileLocation == "" {
		return "."
	}
	return d.DockerfileLocation
}
