package actions

import (
	"os"
	"strings"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	"github.com/arschles/gci/util/docker"
	"github.com/codegangsta/cli"
)

// Test is the CLI handler for 'gci test'
func Test(c *cli.Context) {
	cfg := config.ReadOrDie(c.String(FlagConfigFile))
	paths := PathsOrDie()

	dockerClient := docker.ClientOrDie()
	img, err := docker.ParseImageFromName(docker.GolangImage)
	if err != nil {
		log.Err("error parsing docker image %s (%s)", docker.GolangImage, err)
		os.Exit(1)
	}

	stdOutCh := make(chan docker.Log)
	stdErrCh := make(chan docker.Log)
	exitCodeCh := make(chan int)
	errCh := make(chan error)

	cmd := []string{"go", "test"}
	testPaths := cfg.Test.GetPaths()
	for _, path := range testPaths {
		cmd = append(cmd, path)
	}

	go docker.Run(
		dockerClient,
		img,
		"build",
		paths.CWD,
		docker.ContainerGopath(paths.PackageName),
		strings.Join(cmd, " "),
		cfg.Build.Env,
		stdOutCh,
		stdErrCh,
		exitCodeCh,
		errCh,
	)

	for {
		select {
		case l := <-stdOutCh:
			log.Info("%s", l)
		case l := <-stdErrCh:
			log.Warn("%s", l)
		case err := <-errCh:
			log.Err("%s", err)
			return
		case i := <-exitCodeCh:
			log.Info("exited with code %d", i)
			return
		}
	}

}
