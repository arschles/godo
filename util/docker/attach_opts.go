package docker

import (
	"io"

	docker "github.com/fsouza/go-dockerclient"
)

// AttachContainerOpts returns docker.AttachToContainerOptions with output and error streams turned on
// as well as logs. the returned io.Reader will output both stdout and stderr
func AttachToContainerOpts(containerID string, stdout io.Writer, stderr io.Writer) docker.AttachToContainerOptions {
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
