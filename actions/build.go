package actions

import (
	"io"
	"os"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	dockutil "github.com/arschles/gci/util/docker"
	dockbuild "github.com/arschles/gci/util/docker/build"
	"github.com/codegangsta/cli"
)

// Build is the CLI handler for 'gci build'
func Build(c *cli.Context) {
	cfg := config.ReadOrDie(c.String(FlagConfigFile))
	paths := PathsOrDie()

	dockerClient := dockutil.ClientOrDie()
	imgName := dockutil.ImageName()

	if err := dockutil.EnsureImage(dockerClient, imgName, func() (io.Writer, error) {
		log.Info("Pulling image %s before building", imgName)
		return os.Stdout, nil
	}); err != nil {
		log.Err("Error pulling image %s", imgName)
		os.Exit(1)
	}

	logsCh := make(chan dockbuild.Log)
	resultCh := make(chan int)
	errCh := make(chan error)
	go dockutil.Build(
		dockerClient,
		imgName,
		paths.CWD,
		paths.CWD,
		paths.PackageName,
		"/go",
		cfg,
		logsCh,
		resultCh,
		errCh,
	)

	for {
		select {
		case l := <-logsCh:
			log.Info(l.Message())
		case code := <-resultCh:
			if code == 0 {
				log.Info("Success")
				os.Exit(0)
			} else {
				log.Err("Build failed with code %d", code)
				os.Exit(code)
			}
		case err := <-errCh:
			log.Err("Build failed (%s)", err)
			os.Exit(1)
		}
	}
}
