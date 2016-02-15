package docker

import (
	"github.com/arschles/gci/config"
	docker "github.com/fsouza/go-dockerclient"
)

func Build(dockerCl *docker.Client, rootDir string, outDir string, cfg *config.File) (BuildResult, error) {
	return BuildResult{}, nil
}
