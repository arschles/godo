package actions

import (
	"fmt"
	"os"

	"github.com/arschles/godo/config"
	"github.com/arschles/godo/log"
	"github.com/arschles/godo/util/docker"
	"github.com/codegangsta/cli"
)

const (
	// ListCustomFlag is flag used by godo to list all custom dependencies
	ListCustomFlag = "list"
)

// Custom is the CLI action for 'godo custom ...' commands
func Custom(c *cli.Context) {
	cfg := config.ReadOrDie(c.String(FlagConfigFile))
	if c.Bool(ListCustomFlag) {
		for _, customTarget := range cfg.Custom {
			log.Info("'%s' - %s", customTarget.Name, customTarget.Description)
		}
		return
	}
	if len(c.Args()) < 1 || c.Args()[0] == "" {
		log.Err("you must call this command as 'godo custom <target>'")
		os.Exit(1)
	}
	customName := c.Args()[0]
	customMap := make(map[string]config.CustomTarget)
	for _, customTarget := range cfg.Custom {
		customMap[customTarget.Name] = customTarget
	}
	target, ok := customMap[customName]
	if !ok {
		log.Err("no custom target '%s' found", customName)
		os.Exit(1)
	}

	dockerImage, err := docker.ParseImageFromName(fmt.Sprintf("%s:%s", target.ImageName, target.ImageTag))
	if err != nil {
		log.Err("invalid image name %s:%s (%s)", target.ImageName, target.ImageTag)
		os.Exit(1)
	}

	dockerCl := docker.ClientOrDie()
	paths := PathsOrDie()
	log.Info("executing %s in a %s:%s container", customName, target.ImageName, target.ImageTag)

	rmContainerCh := make(chan func())
	stdOutCh := make(chan docker.Log)
	stdErrCh := make(chan docker.Log)
	exitCodeCh := make(chan int)
	errCh := make(chan error)
	go docker.Run(
		dockerCl,
		dockerImage,
		customName,
		paths.CWD,
		target.MountTarget,
		target.Command,
		target.Envs,
		rmContainerCh,
		stdOutCh,
		stdErrCh,
		exitCodeCh,
		errCh,
	)

	for {
		select {
		case fn := <-rmContainerCh:
			defer fn()
		case l := <-stdOutCh:
			log.Info("%s", l)
		case l := <-stdErrCh:
			log.Warn("%s", l)
		case i := <-exitCodeCh:
			log.Info("exited with code %d", i)
			return
		case err := <-errCh:
			log.Err("%s", err)
			return
		}
	}
}
