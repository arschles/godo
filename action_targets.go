package main

import (
	"github.com/arschles/gci/build"
	"github.com/arschles/gci/log"
	"github.com/codegangsta/cli"
)

func targets(c *cli.Context) {
	buildFile := build.GetFileOrDie(c.GlobalString(flagFile))
	for _, pipeline := range buildFile.Pipelines {
		descr := pipeline.Description
		if descr == "" {
			descr = "[no description]"
		}
		log.Msg("%s - %s", pipeline.Name, descr)
	}
}
