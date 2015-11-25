package main

import (
	"io/ioutil"
	"os"

	"github.com/arschles/gocons/log"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
)

const defaultConsFileName = "gocons.yaml"

func getConsfile() (*Consfile, error) {
	b, err := ioutil.ReadFile(defaultConsFileName)
	if err != nil {
		return nil, err
	}
	f := &Consfile{}
	if err := yaml.Unmarshal(b, f); err != nil {
		return nil, err
	}
	return f, nil
}

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
			Name:        "build",
			Aliases:     []string{"b"},
			Usage:       "build the project",
			Description: "Build code according to the 'build:' directive in the consfile.",
			Action:      build,
		},
		{
			Name:        "bootstrap",
			Usage:       "bootstrap the project",
			Description: "Bootstrap the project. Generally, you'll only have to do this once after you first clone the project.",
			Action:      bootstrap,
		},
		{
			Name:        "other",
			Aliases:     []string{"oth"},
			Usage:       "execute a command listed under 'other'",
			Description: "This command will run a command specified under 'other'",
			Action:      other,
		},
		{
			Name:        "install",
			Aliases:     []string{"i"},
			Usage:       "install gocons",
			Description: "Install gocons using 'go install'",
			Action:      install,
		},
	}

	app.Before = func(c *cli.Context) error {
		log.IsDebugging = c.Bool("debug")
		return nil
	}

	app.Run(os.Args)
}
