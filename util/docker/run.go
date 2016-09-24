package docker

import (
	"io"
	"os"
	"strings"

	"github.com/arschles/godo/log"
	docker "github.com/fsouza/go-dockerclient"
)

// Run runs cmd in the given image using the docker client cl. It mounts cwd into containerMount in the running container and sends on the following channels:
//
// - rmContainerCh: a function closure that the receiver should call, after they receive on errCh or exitCodeCh, to remove the container. this is commonly done with a 'defer'
// - stdOut: all logs from STDOUT in the container. this may never receive
// - stdErr: all logs from STDERR in the container. this may never receive
// - exitCodeCh: the exit code of the container
// - errCh: any error in setting up or running the container. if errCh receives, exitCodeCh may not receive
func Run(
	cl *docker.Client,
	image *Image,
	taskName,
	cwd,
	containerMount,
	cmd string,
	env []string,
	rmContainerCh chan<- func(),
	stdOut chan<- Log,
	stdErr chan<- Log,
	exitCodeCh chan<- int,
	errCh chan<- error,
) {

	mounts := []docker.Mount{
		{Name: "pwd", Source: cwd, Destination: containerMount, Mode: "rxw"},
	}
	cmdSpl := strings.Split(cmd, " ")

	containerName := NewContainerName(taskName, cwd)
	createContainerOpts, hostConfig := CreateAndStartContainerOpts(image.String(), containerName, cmdSpl, env, mounts, containerMount)
	if err := EnsureImage(cl, image.String(), func() (io.Writer, error) {
		return os.Stdout, nil
	}); err != nil {
		errCh <- err
		return
	}

	container, err := cl.CreateContainer(createContainerOpts)
	if err != nil {
		errCh <- err
	}

	rmContainerCh <- func() {
		if err := cl.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID, Force: true}); err != nil {
			log.Warn("Error removing container %s (%s)", container.ID, err)
		}
	}

	log.Debug(CmdStr(createContainerOpts, hostConfig))

	attachOpts := AttachToContainerOpts(container.ID, NewChanWriter(stdOut), NewChanWriter(stdErr))
	// attach before the container starts, so we get all the logs etc...
	go AttachAndWait(cl, container.ID, attachOpts, exitCodeCh, errCh)

	if startErr := cl.StartContainer(container.ID, &hostConfig); startErr != nil {
		errCh <- err
		return
	}
}
