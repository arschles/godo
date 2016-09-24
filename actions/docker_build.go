package actions

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arschles/godo/config"
	"github.com/arschles/godo/log"
	dockutil "github.com/arschles/godo/util/docker"
	"github.com/codegangsta/cli"
	docker "github.com/fsouza/go-dockerclient"
)

// DockerBuild is the CLI action for 'godo docker-build'
func DockerBuild(c *cli.Context) {
	dockerClient := dockutil.ClientOrDie()

	cfg := config.ReadOrDie(c.String(FlagConfigFile))
	if cfg.Docker.ImageName == "" {
		log.Err("Docker image name was empty")
		os.Exit(1)
	}

	dockerfileLocation := cfg.Docker.Build.GetDockerfileLocation()
	dockerfileBytes, err := ioutil.ReadFile(dockerfileLocation)
	if err != nil {
		log.Err("Reading Dockerfile %s [%s]", dockerfileLocation, err)
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

	buildCtx, err := filepath.Abs(cfg.Docker.Build.Context.GetDirectory())
	if err != nil {
		log.Err("Invalid Docker build context %s [%s]", cfg.Docker.Build.Context.GetDirectory(), err)
		os.Exit(1)
	}

	skipSet := cfg.Docker.Build.Context.GetSkips()
	err = tarDir(buildCtx, tr, func(path string, fi os.FileInfo) bool {
		if _, ok := skipSet[path]; ok {
			return true
		}
		if strings.Contains(path, ".git") {
			return true
		}
		if fi.Name() == "Dockerfile" {
			return true
		}
		return false
	})
	if err != nil {
		log.Err("Archiving the build context directory %s [%s]", buildCtx, err)
		os.Exit(1)
	}

	if err := tr.Close(); err != nil {
		log.Err("Closing the build context archive preparing to send it to the Docker daemon [%s]", err)
		os.Exit(1)
	}

	opts := docker.BuildImageOptions{
		Name:           fmt.Sprintf("%s:%s", cfg.Docker.ImageName, cfg.Docker.GetTag()),
		InputStream:    buf,
		Dockerfile:     "Dockerfile",
		OutputStream:   os.Stdout,
		RmTmpContainer: true,
		Pull:           true,
	}
	if err := dockerClient.BuildImage(opts); err != nil {
		log.Err("Building image %s [%s]", cfg.Docker.ImageName, err)
		os.Exit(1)
	}
	log.Info("Successfully built Docker image %s", opts.Name)
}
