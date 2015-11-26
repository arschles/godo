package main

import (
	"github.com/arschles/gocons/log"
	"github.com/codegangsta/cli"
)

func targets(c *cli.Context) {
	consfile := getConsfileOrDie()
	for _, target := range consfile.Targets {
		descr := target.Description
		if descr == "" {
			descr = "[no description]"
		}
		log.Msg("%s - %s", target.Name, descr)
	}
}
