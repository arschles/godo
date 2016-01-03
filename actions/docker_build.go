package actions

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/dockutil"
	"github.com/arschles/gci/log"
	"github.com/codegangsta/cli"
	docker "github.com/fsouza/go-dockerclient"
)

func DockerBuild(c *cli.Context) {
	dockerClient := dockutil.ClientOrDie()

	cfg := config.ReadOrDie(c.String(FlagConfigFile))
	if cfg.DockerBuild.ImageName == "" {
		log.Err("Docker image name was empty")
		os.Exit(1)
	}

	dockerfileBytes, err := ioutil.ReadFile(cfg.DockerBuild.GetDockerfileLocation())
	if err != nil {
		log.Err("Reading Dockerfile %s [%s]", cfg.DockerBuild.GetDockerfileLocation(), err)
		os.Exit(1)
	}

	opts := docker.BuildImageOptions{
		Name:           fmt.Sprintf("%s:%s", cfg.DockerBuild.ImageName, cfg.DockerBuild.GetTag()),
		Dockerfile:     string(dockerfileBytes),
		OutputStream:   os.Stdout,
		RmTmpContainer: true,
		Remote:         fmt.Sprintf("https://%s", cfg.DockerBuild.ImageName),
	}
	if err := dockerClient.BuildImage(opts); err != nil {
		log.Err("Building image %s [%s]", cfg.DockerBuild.ImageName, err)
		os.Exit(1)
	}
}
