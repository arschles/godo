package actions

import (
	"os"
	"strings"

	"github.com/arschles/godo/config"
	"github.com/arschles/godo/log"
	dockutil "github.com/arschles/godo/util/docker"
	"github.com/codegangsta/cli"
	docker "github.com/fsouza/go-dockerclient"
)

// DockerPush is the cli action for 'godo docker-push'
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

	authFileLoc := cfg.Docker.Push.GetAuthFileLocation()
	authFile, err := os.Open(authFileLoc)
	if err != nil {
		log.Err("Reading Docker auth file %s [%s]", authFileLoc, err)
		os.Exit(1)
	}
	defer func() {
		if err := authFile.Close(); err != nil {
			log.Err("Closing Docker auth file %s [%s]", authFileLoc, err)
		}
	}()

	auths, err := docker.NewAuthConfigurations(authFile)
	if err != nil {
		log.Err("Parsing auth file %s [%s]", authFileLoc, err)
		os.Exit(1)
	}

	registry := "https://index.docker.io/v1/"
	spl := strings.Split(pio.Name, "/")
	if len(spl) == 3 {
		registry = spl[0]
	}

	auth, ok := auths.Configs[registry]
	if !ok {
		log.Err("Registry %s in your image %s is not in auth file %s ", registry, pio.Name, authFileLoc)
		os.Exit(1)
	}

	if err := dockerClient.PushImage(pio, auth); err != nil {
		log.Err("Pushing Docker image %s [%s]", pio.Name, err)
		os.Exit(1)
	}
	log.Info("Successfully pushed Docker image %s:%s", pio.Name, pio.Tag)
}
