package actions

import (
	"fmt"
	"os"
	"sync"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	docker "github.com/fsouza/go-dockerclient"
)

type dockerPushOneErr struct {
	err     error
	imgPush config.ImagePush
}

func (e dockerPushOneErr) Error() string {
	return fmt.Sprintf("error pushing %s:%s (%s)", e.imgPush.Name, e.imgPush.Tag, e.err)
}

func dockerPushOne(cl *docker.Client, authCfgs *docker.AuthConfigurations, push config.ImagePush, errCh chan<- error, wg *sync.WaitGroup) {

	defer wg.Done()
	pio := docker.PushImageOptions{
		Name:         push.Name,
		Tag:          push.GetTag(),
		OutputStream: os.Stdout,
	}

	registry := push.GetRegistry()
	authCfg, ok := authCfgs.Configs[registry]
	if !ok {
		log.Debug("no auth config for registry %s, using no auth", registry)
		authCfg = docker.AuthConfiguration{}
	}

	if err := cl.PushImage(pio, authCfg); err != nil {
		log.Err("Pushing Docker image %s [%s]", pio.Name, err)
		errCh <- err
	}
	log.Info("Successfully pushed Docker image %s:%s", pio.Name, pio.Tag)

}
