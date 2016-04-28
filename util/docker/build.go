package docker

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/util/docker/build"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pborman/uuid"
)

const (
	containerOutDir = "/gobin"
)

func goxOutputTpl(binPath string) string {
	return fmt.Sprintf("%s_{{.OS}}_{{.Arch}}", binPath)
}

// ImageName returns the image name to use, given whether we're trying to cross-compile or not
func ImageName(crossCompile bool) string {
	if crossCompile {
		return GoxImage
	}
	return GolangImage
}

func command(crossCompile bool, binaryPath string) []string {
	if crossCompile {
		return []string{"gox", "-output", goxOutputTpl(binaryPath)}
	}
	return []string{"go", "build", "-o", binaryPath}
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
	logsCh chan<- build.Log,
	resultCh chan<- int,
	errCh chan<- error) {

	projName := filepath.Base(rootDir)
	containerName := fmt.Sprintf("gci-build-%s-%s", projName, uuid.New())
	logsCh <- build.LogFromString("Creating container %s to build %s", containerName, packageName)

	binaryName := cfg.Build.GetOutputBinary(projName)
	cmd := command(cfg.Build.CrossCompile, fmt.Sprintf("%s/%s", containerOutDir, binaryName))
	env := cfg.Build.Env

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
		containerGoPath,
		containerWorkDir,
	)
	container, err := dockerCl.CreateContainer(createContainerOpts)
	if err != nil {
		errCh <- fmt.Errorf("error creating container (%s)", err)
		return
	}

	logsCh <- build.LogFromString(CmdStr(createContainerOpts, hostConfig))

	if err := dockerCl.StartContainer(container.ID, &hostConfig); err != nil {
		errCh <- fmt.Errorf("error starting container (%s)", err)
		return
	}

	defer func() {
		if err := dockerCl.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID, Force: true}); err != nil {
			log.Printf("Error removing build container %s (%s)", container.ID, err)
		}
	}()

	stdOut := build.NewChanWriter(logsCh)
	stdErr := build.NewChanWriter(logsCh)
	attachOpts := AttachToContainerOpts(container.ID, stdOut, stdErr)
	waitCodeCh, waitErrCh, err := AttachAndWait(dockerCl, container.ID, attachOpts)

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
