package main

import (
	"github.com/arschles/canta/config"
	"github.com/arschles/canta/log"
	"github.com/codegangsta/cli"
)

func targets(c *cli.Context) {
	buildFile := config.GetFileOrDie(c.GlobalString(flagFile))
	for _, target := range buildFile.Targets {
		descr := target.Description
		if descr == "" {
			descr = "[no description]"
		}
		log.Msg("%s - %s", target.Name, descr)
	}
}
