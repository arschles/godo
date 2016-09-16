package docker

import (
	docker "github.com/fsouza/go-dockerclient"
)

// AttachAndWait attaches to and waits for the container with the given ID using the given client, according to the given options. This function should be started in a goroutine. exitCodeCh will receive if the container completed execution. errCh may receive if the attach couldn't be completed or the container didn't complete properly. if errCh receives, exitCodeCh will not receive
func AttachAndWait(
	dockerClient *docker.Client,
	containerID string,
	attachOpts docker.AttachToContainerOptions,
	exitCodeCh chan<- int,
	errCh chan<- error,
) {
	if err := dockerClient.AttachToContainer(attachOpts); err != nil {
		errCh <- err
		return
	}

	code, err := dockerClient.WaitContainer(containerID)
	if err != nil {
		errCh <- err
		return
	}
	exitCodeCh <- code
}
