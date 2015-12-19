package main

import (
	"github.com/arschles/gci/build"
	"github.com/arschles/gci/log"
	"github.com/codegangsta/cli"
)

func lint(c *cli.Context) {
	log.Die("TODO")
	bfile := build.GetFileOrDie(c.GlobalString(flagFile))

	if bfile.Version > buildFileVersion {
		log.Err("The build file has a higher version (%d) than this build supports (<= %d)", bfile.Version, buildFileVersion)
	}
	//
	// for i, env := range buildFile.Envs {
	// 	if env.Name == "" {
	// 		log.Err("Environment variable #%d has no name", i)
	// 	} else if env.Val == "" && env.Default == "" {
	// 		log.Err("Environment variable #%d (%s) has neither a value nor a default value", i, env.Name)
	// 	}
	// }
	//
	// names := make(map[string]int)
	//
	// for i, target := range buildFile.Targets {
	// 	if target.Name == "" {
	// 		log.Err("Target #%d has no name", i)
	// 		continue
	// 	}
	// 	if prevTargetNum, ok := names[target.Name]; ok {
	// 		log.Err("Target %d uses a name already taken by target %d", i, prevTargetNum)
	// 	}
	// 	names[target.Name] = i
	// 	if target.Description == "" {
	// 		log.Err("Target '%s' (#%d) has no description", target.Name, i)
	// 	}
	// 	if len(target.Commands) == 0 {
	// 		log.Err("Target '%s' (#%d) has no commands", target.Name, i)
	// 	}
	// }
}
