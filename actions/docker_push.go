package actions

import (
	"os"
	"sync"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	"github.com/arschles/gci/util/docker"
	"github.com/codegangsta/cli"
)

// DockerPush is the cli action for 'gci docker-push'
func DockerPush(c *cli.Context) {
	dockerClient := docker.ClientOrDie()
	cfg := config.ReadOrDie(c.String(FlagConfigFile))

	authFileLoc := cfg.Docker.Push.GetAuthFileLocation()
	authCfgs, closeFn, err := getAuthConfigs(authFileLoc)
	if err != nil {
		log.Err("Getting auth file (%s)", err)
		os.Exit(1)
	}
	defer closeFn()

	errCh := make(chan error)
	var wg sync.WaitGroup
	for _, img := range cfg.Docker.Push.Images {
		wg.Add(1)
		go dockerPushOne(dockerClient, authCfgs, img, errCh, &wg)
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
