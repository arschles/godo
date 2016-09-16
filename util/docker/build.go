package docker

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/arschles/gci/config"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pborman/uuid"
)

const (
	containerOutDir = "/gobin"
)

func goxOutputTpl(binPath string) string {
	return fmt.Sprintf("%s_{{.OS}}_{{.Arch}}", binPath)
}

func command(binaryPath string) []string {
	return []string{"go", "build", "-o", binaryPath}
}

func ContainerGopath(packageName string) string {
	return fmt.Sprintf("/go/src/%s", packageName)
}

// Build runs the build of rootDir inside a Docker container, putting binaries into outDir
func Build(
	dockerCl *docker.Client,
	imgName,
	rootDir,
	outDir,
	packageName,
	containerGoPath string,
	cfg *config.File,
	logsCh chan<- Log,
	resultCh chan<- int,
	errCh chan<- error) {

	projName := filepath.Base(rootDir)
	containerName := fmt.Sprintf("gci-build-%s-%s", projName, uuid.New())
	logsCh <- LogFromString("Creating container %s to build %s", containerName, packageName)

	binaryName := cfg.Build.GetOutputBinary(projName)
	cmd := command(fmt.Sprintf("%s/%s", containerOutDir, binaryName))
	env := cfg.Build.Env
	env = append(env, "GOPATH="+containerGoPath)

	containerWorkDir := fmt.Sprintf("%s/src/%s", containerGoPath, packageName)

	mounts := []docker.Mount{
		{
			Name:        "source_dir",
			Source:      rootDir,
			Destination: containerWorkDir,
			Mode:        "r",
		},
		{
			Name:        "dest_dir",
			Source:      outDir,
			Destination: containerOutDir,
			Mode:        "w",
		},
	}
	createContainerOpts, hostConfig := CreateAndStartContainerOpts(
		imgName,
		containerName,
		cmd,
		env,
		mounts,
		containerWorkDir,
	)
	container, err := dockerCl.CreateContainer(createContainerOpts)
	if err != nil {
		errCh <- fmt.Errorf("error creating container (%s)", err)
		return
	}

	logsCh <- LogFromString(CmdStr(createContainerOpts, hostConfig))

	if err := dockerCl.StartContainer(container.ID, &hostConfig); err != nil {
		errCh <- fmt.Errorf("error starting container (%s)", err)
		return
	}

	defer func() {
		if err := dockerCl.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID, Force: true}); err != nil {
			log.Printf("Error removing build container %s (%s)", container.ID, err)
		}
	}()

	stdOut := NewChanWriter(logsCh)
	stdErr := NewChanWriter(logsCh)
	attachOpts := AttachToContainerOpts(container.ID, stdOut, stdErr)
	waitCodeCh := make(chan int)
	waitErrCh := make(chan error)
	go AttachAndWait(dockerCl, container.ID, attachOpts, waitCodeCh, waitErrCh)

	if err != nil {
		errCh <- fmt.Errorf("error attaching to the build container (%s)", err)
		return
	}

	select {
	case err := <-waitErrCh:
		errCh <- fmt.Errorf("error waiting for the build container to finish (%s)", err)
		return
	case code := <-waitCodeCh:
		resultCh <- code
	}
}
