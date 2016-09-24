package config

import (
	"path/filepath"

	"github.com/arschles/godo/util"
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
	DockerfileLocation string             `yaml:"dockerfile-location"`
	Context            DockerBuildContext `yaml:"context"`
}

func (d DockerBuild) GetDockerfileLocation() string {
	if d.DockerfileLocation == "" {
		return "./Dockerfile"
	}
	return d.DockerfileLocation
}

type DockerBuildContext struct {
	Directory string   `yaml:"dir"`
	Skips     []string `yaml:"skip"`
}

func (d DockerBuildContext) GetDirectory() string {
	if d.Directory == "" {
		return "."
	}
	return d.Directory
}

func (d DockerBuildContext) GetSkips() map[string]struct{} {
	ret := make(map[string]struct{})
	for _, sk := range d.Skips {
		ret[sk] = struct{}{}
	}
	return ret
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
