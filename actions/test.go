package actions

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/dockutil"
	"github.com/arschles/gci/log"
	"github.com/codegangsta/cli"
	"github.com/pborman/uuid"
)

func Test(c *cli.Context) {
	cfg := config.ReadOrDie(c.String(FlagConfigFile))
	paths := pathsOrDie()
	dockerClient := dockutil.ClientOrDie()
	projName := filepath.Base(paths.cwd)
	name := fmt.Sprintf("gci-test-%s-%s", projName, uuid.New())
	cmd := []string{"go", "test"}
	for _, path := range cfg.Test.GetPaths() {
		cmd = append(cmd, path)
	}

	createContainerOpts, hostConfig := dockutil.CreateAndStartContainerOpts(
		dockutil.GolangImage,
		name,
		cmd,
		cfg.Test.Env,
		paths.gopath,
		paths.pkg,
	)
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
