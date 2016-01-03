package actions

import (
	"os"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/dockutil"
	"github.com/arschles/gci/log"
	"github.com/codegangsta/cli"
	docker "github.com/fsouza/go-dockerclient"
)

func DockerPush(c *cli.Context) {
	dockerClient := dockutil.ClientOrDie()
	cfg := config.ReadOrDie(c.String(FlagConfigFile))
	if cfg.Docker.ImageName == "" {
		log.Err("Docker image name was empty")
		os.Exit(1)
	}

	pio := docker.PushImageOptions{
		Name:         cfg.Docker.ImageName,
		Tag:          cfg.Docker.GetTag(),
		OutputStream: os.Stdout,
	}

	// TODO: support auth (https://github.com/arschles/gci/issues/16)
	if err := dockerClient.PushImage(pio, docker.AuthConfiguration{}); err != nil {
		log.Err("Pushing Docker image %s [%s]", pio.Name, err)
		os.Exit(1)
	}
	log.Info("Successfully pushed Docker image %s:%s", pio.Name, pio.Tag)
}
