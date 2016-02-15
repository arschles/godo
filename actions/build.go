package actions

import (
	"fmt"
	"os"
	"path/filepath"

	"code.google.com/p/go-uuid/uuid"
	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	dockutil "github.com/arschles/gci/util/docker"
	"github.com/codegangsta/cli"
)

func Build(c *cli.Context) {
	cfg := config.ReadOrDie(c.String(FlagConfigFile))
	paths := pathsOrDie()
	projName := filepath.Base(paths.cwd)
	binary := cfg.Build.GetOutputBinary(projName)

	dockerClient := dockutil.ClientOrDie()

	name := fmt.Sprintf("gci-build-%s-%s", projName, uuid.New())
	cmd := []string{"go", "build", "-o", binary}
	imgName := dockutil.GolangImage
	if cfg.Build.CrossCompile {
		imgName = dockutil.GoxImage
		cmd = []string{"gox"}
	}
	env := cfg.Build.Env
	createContainerOpts, hostConfig := dockutil.CreateAndStartContainerOpts(imgName, name, cmd, env, paths.gopath, paths.pkg)
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
		log.Err("Attaching to the build container [%s]", err)
		os.Exit(1)
	}

	select {
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
