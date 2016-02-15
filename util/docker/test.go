package docker

import (
	"fmt"
	"log"
	"path/filepath"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/pborman/uuid"
)

func testCommand(packages []string) []string {
	ret := make([]string, len(packages)+2)
	ret[0], ret[1] = "go", "test"
	for i, pkg := range packages {
		ret[i+2] = pkg
	}
	return ret
}

// Test runs tests of rootDir inside a Docker container
func Test(
	dockerCl *docker.Client,
	rootDir,
	packageName,
	containerGoPath string,
	packages,
	env []string,
	logsCh chan<- string,
	resultCh chan<- int,
	errCh chan<- error) {

	projName := filepath.Base(rootDir)
	imgName := GolangImage
	containerName := fmt.Sprintf("gci-build-%s-%s", projName, uuid.New())
	logsCh <- fmt.Sprintf("Creating container %s to build %s", containerName, packageName)

	cmd := testCommand(packages)
	containerWorkDir := fmt.Sprintf("%s/src/%s", containerGoPath, packageName)

	mounts := []docker.Mount{
		{
			Name:        "source_dir",
			Source:      rootDir,
			Destination: containerWorkDir,
			Mode:        "r",
		},
	}
	createContainerOpts, hostConfig := CreateAndStartContainerOpts(
		imgName,
		containerName,
		cmd,
		env,
		mounts,
		containerGoPath,
		containerWorkDir,
	)
	container, err := dockerCl.CreateContainer(createContainerOpts)
	if err != nil {
		errCh <- fmt.Errorf("error creating container (%s)", err)
		return
	}

	logsCh <- CmdStr(createContainerOpts, hostConfig)

	if err := dockerCl.StartContainer(container.ID, &hostConfig); err != nil {
		errCh <- fmt.Errorf("error starting container (%s)", err)
		return
	}

	defer func() {
		if err := dockerCl.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID, Force: true}); err != nil {
			log.Printf("Error removing test container %s (%s)", container.ID, err)
		}
	}()

	stdOut := NewChanWriter(logsCh)
	stdErr := NewChanWriter(logsCh)
	attachOpts := AttachToContainerOpts(container.ID, stdOut, stdErr)
	waitCodeCh, waitErrCh, err := AttachAndWait(dockerCl, container.ID, attachOpts)

	if err != nil {
		errCh <- fmt.Errorf("error attaching to the test container (%s)", err)
		return
	}

	select {
	case err := <-waitErrCh:
		errCh <- fmt.Errorf("error waiting for the test container to finish (%s)", err)
		return
	case code := <-waitCodeCh:
		resultCh <- code
	}
}
