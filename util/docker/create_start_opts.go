package docker

import (
	"fmt"
	"path/filepath"

	docker "github.com/fsouza/go-dockerclient"
)

var (
	defaultEnv = []string{"GO15VENDOREXPERIMENT=1", "CGO_ENABLED=0", "GOPATH=/go"}
)

func CreateAndStartContainerOpts(
	imageName,
	containerName string,
	cmd []string,
	env []string,
	gopath,
	packagePath string,
) (docker.CreateContainerOptions, docker.HostConfig) {

	absPwd := filepath.Join(gopath, "src", packagePath)

	mount := docker.Mount{
		Name:        "pwd",
		Source:      absPwd,
		Destination: filepath.Join(containerGopath, "src", packagePath),
		Mode:        "rx",
	}
	createOpts := docker.CreateContainerOptions{
		Name: containerName,
		Config: &docker.Config{
			Env:   append(defaultEnv, env...),
			Cmd:   cmd,
			Image: imageName,
			Volumes: map[string]struct{}{
				absPwd: struct{}{},
			},
			Mounts:     []docker.Mount{mount},
			WorkingDir: mount.Destination,
		},
		HostConfig: &docker.HostConfig{},
	}
	hostConfig := docker.HostConfig{
		Binds: []string{fmt.Sprintf("%s:%s", mount.Source, mount.Destination)},
	}
	return createOpts, hostConfig
}
