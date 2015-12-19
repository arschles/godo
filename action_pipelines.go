package main

import (
	"fmt"

	"github.com/arschles/gci/build"
	"github.com/arschles/gci/log"
	"github.com/codegangsta/cli"
)

func pipelines(c *cli.Context) {
	bfile := build.GetFileOrDie(c.GlobalString(flagFile))
	varMap := bfile.GetVarMap()
	for _, pipeline := range bfile.Pipelines {
		descr, err := pipeline.RenderDescription(varMap)
		if err != nil {
			descr = fmt.Sprintf("error rendering description [%s]", err)
		} else if descr == "" {
			descr = "[no description]"
		}
		log.Msg("%s - %s", pipeline.Name, descr)
	}
}
