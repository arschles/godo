package actions

import (
	"bufio"
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

func createContainerOpts(pathAboveSrc, pathUnderSrc string) docker.CreateContainerOptions {
	absCwd := filepath.Join(pathAboveSrc, "src", pathUnderSrc)
	projName := filepath.Base(pathUnderSrc)

	return docker.CreateContainerOptions{
		Name: fmt.Sprintf("gci-build-%s-%s", projName, uuid.New()),
		Config: &docker.Config{
			Env:   []string{"GO15VENDOREXPERIMENT=1", "CGO_ENABLED=0", "GOPATH=/go"},
			Cmd:   []string{"go", "build"},
			Image: golangImage,
			Volumes: map[string]struct{}{
				absPwd: struct{}{},
			},
			Mounts: []docker.Mount{
				docker.Mount{
					Name:        "pwd",
					Source:      absPwd,
					Destination: fmt.Sprintf("%s/src/%s", containerGopath, pathUnderSrc),
					Mode:        "rx",
				},
			},
		},
		HostConfig: &docker.HostConfig{},
	}
}

// attachContainerOpts returns docker.AttachToContainerOptions with output and error streams turned on
// as well as logs. the returned io.Reader will output both stdout and stderr
func attachToContainerOpts(containerID string) (docker.AttachToContainerOptions, io.Reader) {
	r, w := io.Pipe()
	// var stdoutBuf, stderrBuf bytes.Buffer
	opts := docker.AttachToContainerOptions{
		Container:    containerID,
		OutputStream: w,
		ErrorStream:  w,
		Logs:         true,
		Stream:       true,
		Stdout:       true,
		Stderr:       true,
	}

	return opts, r
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

	srcSpl := strings.Split(cwd, "src")
	if len(srcSpl) < 2 {
		log.Err("")
		os.Exit(1)
	}

	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		log.Err("creating new docker client (%s)", err)
		os.Exit(1)
	}

	containerOpts := createContainerOpts(cwd)
	container, err := dockerClient.CreateContainer(containerOpts)
	if err != nil {
		log.Err("creating container [%s]", err)
		os.Exit(1)
	}

	hostConfig := &docker.HostConfig{Binds: []string{fmt.Sprintf("%s:%s", workdir, absPwd)}}
	if err := dockerCl.StartContainer(container.ID, hostConfig); err != nil {
		log.Err("starting container [%s]", err)
		os.Exit(1)
	}

	attachOpts, outputReader := attachToContainerOpts(container.ID)
	errCh := make(chan error)
	go func() {
		if err := dockerCl.AttachToContainer(attachOpts); err != nil {
			errCh <- err
		}
	}()

	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			fmt.Fprintf(w, "%s\n", scanner.Text())
			flusher.Flush()
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(w, "error with scanner in attached container [%s]\n", err)
		}
	}(outputReader)

	code, err := dockerClient.WaitContainer(container.ID)
	if err != nil {
		log.Errf("waiting for container %s [%s]", container.ID, err)
		return
	}

}
