package docker

import (
	"os"

	"github.com/arschles/gci/log"
	docker "github.com/fsouza/go-dockerclient"
)

func ClientOrDie() *docker.Client {
	cl, err := docker.NewClientFromEnv()
	if err != nil {
		log.Err("creating new docker client (%s)", err)
		os.Exit(1)
	}
	return cl
}
