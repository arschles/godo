package dockutil

import (
	"fmt"
	"path/filepath"

	docker "github.com/fsouza/go-dockerclient"
)

func CreateAndStartContainerOpts(name string, cmd []string, gopath, packagePath string) (docker.CreateContainerOptions, docker.HostConfig) {
	absPwd := filepath.Join(gopath, "src", packagePath)

	mount := docker.Mount{
		Name:        "pwd",
		Source:      absPwd,
		Destination: filepath.Join(containerGopath, "src", packagePath),
		Mode:        "rx",
	}
	createOpts := docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Env:   []string{"GO15VENDOREXPERIMENT=1", "CGO_ENABLED=0", "GOPATH=/go"},
			Cmd:   cmd,
			Image: golangImage,
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
