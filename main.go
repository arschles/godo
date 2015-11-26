package main

import (
	"os"

	"github.com/arschles/gocons/log"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gocons"
	app.Usage = "gocons is a Makefile replacement for Go projects"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable verbose debugging output",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "targets",
			Aliases:     []string{"tgt"},
			Usage:       "list all the targets in this projects",
			Description: "This command will print a list of all targets in the build file (including builtins), along with a short description of each",
			Action:      targets,
		},
		{
			Name:        "run",
			Aliases:     []string{"r"},
			Usage:       "run a target",
			Description: "Run a target. You can list all targets by running 'gocons targets'",
			Action:      run,
		},
	}

	app.Before = func(c *cli.Context) error {
		log.IsDebugging = c.Bool("debug")
		return nil
	}

	app.Run(os.Args)
}
