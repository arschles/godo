package actions

import (
	"fmt"
	"os"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	"github.com/arschles/gci/util/docker"
	"github.com/codegangsta/cli"
)

const (
	// ListCustomFlag is flag used by gci to list all custom dependencies
	ListCustomFlag = "list"
)

// Custom is the CLI action for 'gci custom ...' commands
func Custom(c *cli.Context) {
	cfg := config.ReadOrDie(c.String(FlagConfigFile))
	if c.Bool(ListCustomFlag) {
		for _, customTarget := range cfg.Custom {
			log.Info("'%s' - %s", customTarget.Name, customTarget.Description)
		}
		return
	}
	if len(c.Args()) < 1 || c.Args()[0] == "" {
		log.Err("you must call this command as 'gci custom <target>'")
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
	exitCode, runErr := docker.Run(dockerCl, dockerImage, customName, paths.CWD, target.MountTarget, target.Command, target.Envs, os.Stdout)
	if runErr != nil {
		log.Err("running '%s' (%s)", target.Command, runErr)
		os.Exit(1)
	}
	log.Info("'%s' exited with code %d", target.Command, exitCode)
}
