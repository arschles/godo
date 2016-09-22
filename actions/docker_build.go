package actions

import (
	"sync"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	dockutil "github.com/arschles/gci/util/docker"
	"github.com/codegangsta/cli"
)

// DockerBuild is the CLI action for 'gci docker-build'
func DockerBuild(c *cli.Context) {
	dockerClient := dockutil.ClientOrDie()

	cfg := config.ReadOrDie(c.String(FlagConfigFile))

	errCh := make(chan error)
	var wg sync.WaitGroup
	for _, imgBuild := range cfg.Docker.Build.Images {
		wg.Add(1)
		go dockerBuildOne(dockerClient, imgBuild, errCh, &wg)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		log.Err("%s", err)
	}
	log.Info("done")
}
