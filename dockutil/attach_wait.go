package dockutil

import (
	docker "github.com/fsouza/go-dockerclient"
)

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
