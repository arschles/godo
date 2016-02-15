package docker

import (
	"fmt"
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

func imageName(crossCompile bool) string {
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
	rootDir,
	outDir,
	packageName string,
	cfg *config.File,
	logsCh chan<- build.Log,
	resultCh chan<- int,
	errCh chan<- error) {

	projName := filepath.Base(rootDir)
	imgName := imageName(cfg.Build.CrossCompile)
	containerName := fmt.Sprintf("gci-build-%s-%s", projName, uuid.New())
	logsCh <- build.LogFromString("Creating container %s", containerName)

	binaryName := cfg.Build.GetOutputBinary(projName)
	cmd := command(cfg.Build.CrossCompile, fmt.Sprintf("%s/%s", containerOutDir, binaryName))
	env := cfg.Build.Env

	containerWorkDir := fmt.Sprintf("%s/%s", ContainerGoPath, packageName)

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
		ContainerGoPath,
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
