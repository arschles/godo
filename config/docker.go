package config

import (
	"path/filepath"

	"github.com/arschles/gci/util"
)

type Docker struct {
	ImageName string      `yaml:"image-name"`
	Tag       string      `yaml:"tag"`
	Build     DockerBuild `yaml:"build"`
	Push      DockerPush  `yaml:"push"`
}

func (d Docker) GetTag() string {
	if d.Tag == "" {
		return "latest"
	}
	return d.Tag
}

type DockerBuild struct {
	DockerfileLocation string `yaml:"dockerfile-location"`
	ContextPath        string `yaml:"context-path"`
}

func (d DockerBuild) GetContextPath() string {
	if d.ContextPath == "" {
		return "."
	}
	return d.ContextPath
}

func (d DockerBuild) GetDockerfileLocation() string {
	if d.DockerfileLocation == "" {
		return "./Dockerfile"
	}
	return d.DockerfileLocation
}

type DockerPush struct {
	AuthFileLocation string `json:"auth-file-location"`
}

func (d DockerPush) GetAuthFileLocation() string {
	if d.AuthFileLocation == "" {
		return filepath.Join(util.GetHome(), ".docker", "config.json")
	}
	return d.AuthFileLocation
}
