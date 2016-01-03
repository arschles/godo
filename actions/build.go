package actions

import (
	"fmt"
	"os"
	"path/filepath"

	"code.google.com/p/go-uuid/uuid"
	"github.com/arschles/gci/dockutil"
	"github.com/arschles/gci/log"
	"github.com/codegangsta/cli"
)

func Build(c *cli.Context) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Err("GOPATH environment variable not found")
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Err("getting current working dir (%s)", err)
		os.Exit(1)
	}
	projName := filepath.Base(cwd)

	pkgPath, err := packagePath(gopath, cwd)
	if err != nil {
		log.Err("Error detecting package name [%s]", err)
		os.Exit(1)
	}

	dockerClient := dockutil.ClientOrDie()

	name := fmt.Sprintf("gci-build-%s-%s", projName, uuid.New())
	cmd := []string{"go", "build", "-o", projName}
	createContainerOpts, hostConfig := dockutil.CreateAndStartContainerOpts(name, cmd, gopath, pkgPath)
	container, err := dockerClient.CreateContainer(createContainerOpts)
	if err != nil {
		log.Err("creating container [%s]", err)
		os.Exit(1)
	}

	log.Msg(dockutil.CmdStr(createContainerOpts, hostConfig))

	if err := dockerClient.StartContainer(container.ID, &hostConfig); err != nil {
		log.Err("starting container [%s]", err)
		os.Exit(1)
	}

	attachOpts := dockutil.AttachToContainerOpts(container.ID, os.Stdout, os.Stderr)
	attachErrCh := make(chan error)
	go func() {
		if err := dockerClient.AttachToContainer(attachOpts); err != nil {
			attachErrCh <- err
		}
	}()

	waitErrCh := make(chan error)
	waitCodeCh := make(chan int)

	go func() {
		code, err := dockerClient.WaitContainer(container.ID)
		if err != nil {
			waitErrCh <- err
			return
		}
		waitCodeCh <- code
	}()

	select {
	case err := <-attachErrCh:
		log.Err("Attaching to the build container [%s]", err)
		os.Exit(1)
	case err := <-waitErrCh:
		log.Err("Waiting for the build container to finish [%s]", err)
		os.Exit(1)
	case code := <-waitCodeCh:
		if code != 0 {
			log.Err("Build exited %d", code)
			os.Exit(code)
		} else {
			log.Info("Success")
		}
	}
}
