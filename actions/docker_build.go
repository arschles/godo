package actions

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/dockutil"
	"github.com/arschles/gci/log"
	"github.com/codegangsta/cli"
	docker "github.com/fsouza/go-dockerclient"
)

func DockerBuild(c *cli.Context) {
	dockerClient := dockutil.ClientOrDie()

	cfg := config.ReadOrDie(c.String(FlagConfigFile))
	if cfg.Docker.ImageName == "" {
		log.Err("Docker image name was empty")
		os.Exit(1)
	}

	dockerfileBytes, err := ioutil.ReadFile(cfg.Docker.Build.GetDockerfileLocation())
	if err != nil {
		log.Err("Reading Dockerfile %s [%s]", cfg.Docker.Build.GetDockerfileLocation(), err)
		os.Exit(1)
	}

	t := time.Now()
	buf := bytes.NewBuffer(nil)
	tr := tar.NewWriter(buf)
	tr.WriteHeader(&tar.Header{
		Name:       "Dockerfile",
		Size:       int64(len(dockerfileBytes)),
		ModTime:    t,
		AccessTime: t,
		ChangeTime: t,
	})
	tr.Write(dockerfileBytes)

	//TODO: write build context to tar stream as well

	opts := docker.BuildImageOptions{
		Name:           fmt.Sprintf("%s:%s", cfg.Docker.ImageName, cfg.Docker.GetTag()),
		InputStream:    buf,
		OutputStream:   os.Stdout,
		RmTmpContainer: true,
		Pull:           true,
	}
	if err := dockerClient.BuildImage(opts); err != nil {
		log.Err("Building image %s [%s]", cfg.Docker.ImageName, err)
		os.Exit(1)
	}
}
