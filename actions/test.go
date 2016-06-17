package actions

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	dockutil "github.com/arschles/gci/util/docker"
	"github.com/codegangsta/cli"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pborman/uuid"
)

// Test is the CLI action for 'gci test'
func Test(c *cli.Context) {
	cfg := config.ReadOrDie(c.String(FlagConfigFile))
	paths := PathsOrDie()
	dockerClient := dockutil.ClientOrDie()
	projName := filepath.Base(paths.CWD)
	imgName := dockutil.GolangImage
	name := fmt.Sprintf("gci-test-%s-%s", projName, uuid.New())
	cmd := []string{"go", "test"}
	for _, path := range cfg.Test.GetPaths() {
		cmd = append(cmd, path)
	}

	// TODO: don't assume that the GOPATH in the container is this. Somehow the util/docker package needs to specify it
	workDir := containerGoPath + "/src/" + paths.PackageName
	mounts := []docker.Mount{
		{
			Name:        "pwd",
			Source:      paths.CWD,
			Destination: workDir,
			Mode:        "rx",
		},
	}
	cfg.Test.Env = append(cfg.Test.Env, "GOPATH="+containerGoPath)
	createContainerOpts, hostConfig := dockutil.CreateAndStartContainerOpts(
		imgName,
		name,
		cmd,
		cfg.Test.Env,
		mounts,
		workDir,
	)

	if err := dockutil.EnsureImage(dockerClient, imgName, func() (io.Writer, error) {
		log.Info("Pulling image %s before testing", imgName)
		return os.Stdout, nil
	}); err != nil {
		log.Err("Error pulling image %s", imgName)
	}

	container, err := dockerClient.CreateContainer(createContainerOpts)
	if err != nil {
		log.Err("creating container [%s]", err)
		os.Exit(1)
	}

	log.Msg(dockutil.CmdStr(createContainerOpts, hostConfig))

	if startErr := dockerClient.StartContainer(container.ID, &hostConfig); startErr != nil {
		log.Err("starting container [%s]", startErr)
		os.Exit(1)
	}

	attachOpts := dockutil.AttachToContainerOpts(container.ID, os.Stdout, os.Stderr)
	waitCodeCh, waitErrCh, err := dockutil.AttachAndWait(dockerClient, container.ID, attachOpts)
	if err != nil {
		log.Err("Attaching to the test container [%s]", err)
		os.Exit(1)
	}

	select {
	case err := <-waitErrCh:
		log.Err("Waiting for the test container to finish [%s]", err)
		os.Exit(1)
	case code := <-waitCodeCh:
		if code != 0 {
			log.Err("Test exited %d", code)
			os.Exit(code)
		} else {
			log.Info("Success")
		}
	}
}
