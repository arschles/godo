package actions

import (
	"github.com/arschles/gci/files/build"
	"github.com/arschles/gci/log"
	"github.com/codegangsta/cli"
)

func Run(c *cli.Context) {
	bfile := build.GetFileOrDie(c.GlobalString(FlagFile))
	pName := c.Args().Get(0)
	if pName == "" {
		log.Die("no pipeline name given")
	}
	var pipeline *build.Pipeline = nil
	for _, p := range bfile.Pipelines {
		if p.Name == pName {
			pipeline = &p
			break
		}
	}
	if pipeline == nil {
		log.Die("there's no '%s' pipeline in the build file", pName)
	}

	log.Info("executing %s", pipeline.Name)
	for i, step := range pipeline.Steps {
		log.Info("step %d (%s)", i, step.Name)
	}

	log.Info("done")
}
