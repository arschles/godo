package dockutil

import (
	docker "github.com/fsouza/go-dockerclient"
)

// AttachAndWait attaches to, and waits for, the container with the given ID using the given client, according to the given options.
// The first returned channel will receive if there was an error attaching. The second channel returned will receive if there was an error waiting. The 3rd channel will return with the exit code.
func AttachAndWait(dockerClient *docker.Client, containerID string, attachOpts docker.AttachToContainerOptions) (<-chan error, <-chan error, <-chan int) {
	attachErrCh := make(chan error)
	go func() {
		if err := dockerClient.AttachToContainer(attachOpts); err != nil {
			attachErrCh <- err
		}
	}()

	waitErrCh := make(chan error)
	waitCodeCh := make(chan int)

	go func() {
		code, err := dockerClient.WaitContainer(containerID)
		if err != nil {
			waitErrCh <- err
			return
		}
		waitCodeCh <- code
	}()

	return attachErrCh, waitErrCh, waitCodeCh
}
