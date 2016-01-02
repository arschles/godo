package actions

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"code.google.com/p/go-uuid/uuid"
	docker "github.com/fsouza/go-dockerclient"
)

const (
	golangImage     = "golang:1.5.2"
	containerGopath = "/go"
)

func dockerCmd(co docker.CreateContainerOptions, hc docker.HostConfig) string {
	ret := []string{"docker run"}
	for _, env := range co.Config.Env {
		ret = append(ret, fmt.Sprintf("-e %s", env))
	}

	for _, b := range hc.Binds {
		ret = append(ret, fmt.Sprintf("-v %s", b))
	}
	ret = append(ret, fmt.Sprintf("-w %s", co.Config.WorkingDir))
	ret = append(ret, fmt.Sprintf("--name=%s", co.Name))
	ret = append(ret, co.Config.Image)
	ret = append(ret, strings.Join(co.Config.Cmd, " "))

	return strings.Join(ret, " ")
}

func createAndStartContainerOpts(gopath, packagePath string) (docker.CreateContainerOptions, docker.HostConfig) {
	absPwd := filepath.Join(gopath, "src", packagePath)
	projName := filepath.Base(packagePath)

	mount := docker.Mount{
		Name:        "pwd",
		Source:      absPwd,
		Destination: filepath.Join(containerGopath, "src", packagePath),
		Mode:        "rx",
	}
	createOpts := docker.CreateContainerOptions{
		Name: fmt.Sprintf("gci-build-%s-%s", projName, uuid.New()),
		Config: &docker.Config{
			Env:   []string{"GO15VENDOREXPERIMENT=1", "CGO_ENABLED=0", "GOPATH=/go"},
			Cmd:   []string{"go", "build"},
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

// attachContainerOpts returns docker.AttachToContainerOptions with output and error streams turned on
// as well as logs. the returned io.Reader will output both stdout and stderr
func attachToContainerOpts(containerID string, stdout io.Writer, stderr io.Writer) docker.AttachToContainerOptions {
	// var stdoutBuf, stderrBuf bytes.Buffer
	opts := docker.AttachToContainerOptions{
		Container:    containerID,
		OutputStream: stdout,
		ErrorStream:  stderr,
		Logs:         true,
		Stream:       true,
		Stdout:       true,
		Stderr:       true,
	}

	return opts
}
