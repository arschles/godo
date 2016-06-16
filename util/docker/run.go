package docker

import (
	"io"
	"os"
	"strings"

	"github.com/arschles/gci/log"
	docker "github.com/fsouza/go-dockerclient"
)

// Run runs cmd in the given image using the docker client cl. It mounts cwd into containerMount in the running container, and outputs the result of the run to out. After the container finishes running, returns the command's exit code. If there was an error starting or running it, returns 0 and a non-nil error
func Run(cl *docker.Client, image *Image, taskName, cwd, containerMount, cmd string, env []string, out io.Writer) (int, error) {
	mounts := []docker.Mount{
		{Name: "pwd", Source: cwd, Destination: containerMount, Mode: "rx"},
	}
	cmdSpl := strings.Split(cmd, " ")

	containerName := NewContainerName(taskName, cwd)
	createContainerOpts, hostConfig := CreateAndStartContainerOpts(image.String(), containerName, cmdSpl, env, mounts, containerMount)
	if err := EnsureImage(cl, image.String(), func() (io.Writer, error) {
		return os.Stdout, nil
	}); err != nil {
		return 0, err
	}

	container, err := cl.CreateContainer(createContainerOpts)
	if err != nil {
		return 0, err
	}

	log.Msg(CmdStr(createContainerOpts, hostConfig))

	if startErr := cl.StartContainer(container.ID, &hostConfig); startErr != nil {
		return 0, startErr
	}

	attachOpts := AttachToContainerOpts(container.ID, os.Stdout, os.Stderr)
	waitCodeCh, waitErrCh, err := AttachAndWait(cl, container.ID, attachOpts)
	if err != nil {
		return 0, err
	}

	select {
	case err := <-waitErrCh:
		return 0, err
	case code := <-waitCodeCh:
		return code, nil
	}
}
