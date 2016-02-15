package actions

import (
	"os"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	dockutil "github.com/arschles/gci/util/docker"
	"github.com/codegangsta/cli"
)

func Build(c *cli.Context) {
	cfg := config.ReadOrDie(c.String(FlagConfigFile))
	paths := PathsOrDie()

	dockerClient := dockutil.ClientOrDie()

	logsCh := make(chan string)
	resultCh := make(chan int)
	errCh := make(chan error)
	go dockutil.Build(dockerClient, paths.CWD, paths.CWD, paths.PackageName, "/go", cfg, logsCh, resultCh, errCh)

	for {
		select {
		case l := <-logsCh:
			log.Info(l)
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
