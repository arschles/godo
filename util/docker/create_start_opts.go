package docker

import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
)

var (
	defaultEnv = []string{"GO15VENDOREXPERIMENT=1", "GOPATH=/go"}
)

// CreateAndStartContainerOpts creates a container from imageName with the name containerName. The container will execute cmd with the given enviroment variables (env).
// The container will specify volumes as each 'Source' field in mounts, and bind each mount.Source to each mount.Destination. Finally, the container will set GOPATH to the given containerGoPath variable.
func CreateAndStartContainerOpts(
	imageName,
	containerName string,
	cmd []string,
	env []string,
	mounts []docker.Mount,
	workDir string,
) (docker.CreateContainerOptions, docker.HostConfig) {

	vols := make(map[string]struct{})
	for _, mount := range mounts {
		vols[mount.Source] = struct{}{}
	}
	binds := make([]string, len(mounts))
	for i, mount := range mounts {
		binds[i] = fmt.Sprintf("%s:%s", mount.Source, mount.Destination)
	}

	hostConfig := docker.HostConfig{
		Binds: binds,
	}
	createOpts := docker.CreateContainerOptions{
		Name: containerName,
		Config: &docker.Config{
			Env:        append(defaultEnv, env...),
			Cmd:        cmd,
			Image:      imageName,
			Volumes:    vols,
			Mounts:     mounts,
			WorkingDir: workDir,
		},
		HostConfig: &hostConfig,
	}
	return createOpts, hostConfig
}
