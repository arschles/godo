package main

import (
	"os"

	"github.com/arschles/gci/log"
	"github.com/codegangsta/cli"
)

const (
	buildFileVersion = 1
	flagDebug        = "debug"
	flagFile         = "file"
)

func main() {
	app := cli.NewApp()
	app.Name = "gci"
	app.Usage = "gci is a build and CI tool for Go projects"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  flagDebug,
			Usage: "Enable verbose debugging output",
		},
		cli.StringFlag{
			Name:  flagFile,
			Usage: "Specify the build file to use",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "targets",
			Aliases:     []string{"t"},
			Usage:       "list all the targets in this projects",
			Description: "This command will print a list of all targets in the build file (including builtins), along with a short description of each",
			Action:      targets,
		},
		{
			Name:        "run",
			Aliases:     []string{"r"},
			Usage:       "run a target",
			Description: "You can list all target names and descriptions by running 'gci targets'",
			Action:      run,
		},
		{
			Name:        "lint",
			Aliases:     []string{"l"},
			Usage:       "run a linter on the cons file",
			Description: "The linter checks for a malformed cons file or missing information",
			Action:      lint,
		},
	}

	app.Before = func(c *cli.Context) error {
		log.IsDebugging = c.Bool(flagDebug)
		return nil
	}

	app.Run(os.Args)
}
