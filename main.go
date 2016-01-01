package main

import (
	"os"

	"github.com/arschles/gci/actions"
	"github.com/arschles/gci/log"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gci"
	app.Usage = "gci is a build and CI tool for Go projects"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  actions.FlagDebug,
			Usage: "Enable verbose debugging output",
		},
		cli.StringFlag{
			Name:  actions.FlagFile,
			Usage: "Specify the build file to use",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "build",
			Aliases:     []string{"b"},
			Usage:       "Build your project",
			Description: "This command will build your code from the current working directory",
			Action:      actions.Build,
		},
		{
			Name:        "pipelines",
			Aliases:     []string{"p"},
			Usage:       "list all the pipelines in this project",
			Description: "This command will print a list of all of the pipelines defined in the build file",
			Action:      actions.Pipelines,
		},
		{
			Name:        "run",
			Aliases:     []string{"r"},
			Usage:       "run a target",
			Description: "You can list all target names and descriptions by running 'gci targets'",
			Action:      actions.Run,
		},
		{
			Name:        "lint",
			Aliases:     []string{"l"},
			Usage:       "run a linter on the cons file",
			Description: "The linter checks for a malformed cons file or missing information",
			Action:      actions.Lint,
		},
	}

	app.Before = func(c *cli.Context) error {
		log.IsDebugging = c.Bool(actions.FlagDebug)
		return nil
	}

	app.Run(os.Args)
}
