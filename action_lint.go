package main

import (
	"github.com/arschles/canta/log"
	"github.com/codegangsta/cli"
)

func lint(c *cli.Context) {
	consfile, err := getConsfile()
	if err != nil {
		log.Die("The consfile is missing [%s]", err)
	}

	if consfile.Version > consfileVersion {
		log.Err("The consfile has a higher version (%d) than this build of gocons supports (<= %d)", consfile.Version, consfileVersion)
	}

	for i, env := range consfile.Envs {
		if env.Name == "" {
			log.Err("Environment variable #%d has no name", i)
		} else if env.Val == "" && env.Default == "" {
			log.Err("Environment variable #%d (%s) has neither a value nor a default value", i, env.Name)
		}
	}

	names := make(map[string]int)

	for i, target := range consfile.Targets {
		if target.Name == "" {
			log.Err("Target #%d has no name", i)
			continue
		}
		if prevTargetNum, ok := names[target.Name]; ok {
			log.Err("Target %d uses a name already taken by target %d", i, prevTargetNum)
		}
		names[target.Name] = i
		if target.Description == "" {
			log.Err("Target '%s' (#%d) has no description", target.Name, i)
		}
		if len(target.Commands) == 0 {
			log.Err("Target '%s' (#%d) has no commands", target.Name, i)
		}
	}
}
