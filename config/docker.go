package config

import (
	"path/filepath"
	"strings"

	"github.com/arschles/gci/util"
)

type Docker struct {
	Build DockerBuild `yaml:"build"`
	Push  DockerPush  `yaml:"push"`
}

type ImageBuild struct {
	Name    string             `yaml:"name"`
	Tag     string             `yaml:"tag"`
	Context DockerBuildContext `yaml:"context"`
}

func (i ImageBuild) GetTag() string {
	if i.Tag == "" {
		return "latest"
	}
	return i.Tag
}

func (i ImageBuild) GetDockerfileLocation() string {
	if i.Context.DockerfileLocation == "" {
		return "./Dockerfile"
	}
	return i.Context.DockerfileLocation
}

type DockerBuild struct {
	Images []ImageBuild `yaml:"images"`
}

type DockerBuildContext struct {
	Directory          string   `yaml:"dir"`
	DockerfileLocation string   `yaml:"dockerfile-location"`
	Skips              []string `yaml:"skip"`
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

type ImagePush struct {
	Name string `yaml:"name"`
	Tag  string `yaml:"tag"`
}

func (i ImagePush) GetTag() string {
	if i.Tag == "" {
		return "latest"
	}
	return i.Tag
}

func (i ImagePush) GetRegistry() string {
	registry := "https://index.docker.io/v1/"
	spl := strings.Split(i.Name, "/")
	if len(spl) == 3 {
		registry = spl[0]
	}
	return registry
}

type DockerPush struct {
	AuthFileLocation string      `json:"auth-file-location"`
	Images           []ImagePush `yaml:"images"`
}

func (d DockerPush) GetAuthFileLocation() string {
	if d.AuthFileLocation == "" {
		return filepath.Join(util.GetHome(), ".docker", "config.json")
	}
	return d.AuthFileLocation
}
