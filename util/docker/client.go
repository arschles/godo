package docker

import (
	"os"

	"github.com/arschles/godo/log"
	docker "github.com/fsouza/go-dockerclient"
)

// ClientOrDie creates a new Docker client. If one couldn't be created, logs and error and exits with status code 1
func ClientOrDie() *docker.Client {
	cl, err := docker.NewClientFromEnv()
	if err != nil {
		log.Err("creating new docker client (%s)", err)
		os.Exit(1)
	}
	return cl
}
