package actions

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"code.google.com/p/go-uuid/uuid"
	"github.com/arschles/gci/log"
	"github.com/codegangsta/cli"
	docker "github.com/fsouza/go-dockerclient"
)

const (
	golangImage     = "golang:1.5.2"
	containerGopath = "/go"
)

func createAndStartContainerOpts(gopath, packagePath string) (docker.CreateContainerOptions, docker.HostConfig) {
	absPwd := filepath.Join(gopath, "src", packagePath)
	projName := filepath.Base(packagePath)

	mount := docker.Mount{
		Name:        "pwd",
		Source:      absPwd,
		Destination: fmt.Sprintf("%s/src/%s", containerGopath, packagePath),
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
			Mounts: []docker.Mount{mount},
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

func Build(c *cli.Context) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Err("GOPATH environment variable not found")
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Err("getting current working dir (%s)", err)
		os.Exit(1)
	}

	pkgPath, err := packagePath(gopath, cwd)
	if err != nil {
		log.Err("Error detecting package name [%s]", err)
		os.Exit(1)
	}

	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		log.Err("creating new docker client (%s)", err)
		os.Exit(1)
	}

	createContainerOpts, hostConfig := createAndStartContainerOpts(gopath, pkgPath)
	container, err := dockerClient.CreateContainer(createContainerOpts)
	if err != nil {
		log.Err("creating container [%s]", err)
		os.Exit(1)
	}

	if err := dockerClient.StartContainer(container.ID, &hostConfig); err != nil {
		log.Err("starting container [%s]", err)
		os.Exit(1)
	}

	attachOpts := attachToContainerOpts(container.ID, os.Stdout, os.Stderr)
	attachErrCh := make(chan error)
	go func() {
		if err := dockerClient.AttachToContainer(attachOpts); err != nil {
			attachErrCh <- err
		}
	}()

	waitErrCh := make(chan error)
	waitCodeCh := make(chan int)

	go func() {
		code, err := dockerClient.WaitContainer(container.ID)
		if err != nil {
			waitErrCh <- err
			return
		}
		waitCodeCh <- code
	}()

	select {
	case err := <-attachErrCh:
		log.Err("Attaching to the build container [%s]", err)
		os.Exit(1)
	case err := <-waitErrCh:
		log.Err("Waiting for the build container to finish [%s]", err)
		os.Exit(1)
	case code := <-waitCodeCh:
		if code != 0 {
			log.Err("Build exited %d", code)
			os.Exit(code)
		} else {
			log.Info("Success")
		}
	}
}
