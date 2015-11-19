package main

import (
	"io/ioutil"
	"os"

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
	app.Commands = []cli.Command{
		{
			Name:        "build",
			Aliases:     []string{"bld"},
			Usage:       "build the project",
			Description: "This command will build code according to the 'build:' directive in the consfile.",
			Action:      build,
		},
		{
			Name:        "bootstrap",
			Aliases:     []string{"bts"},
			Usage:       "bootstrap the project",
			Description: "This command will bootstrap the project. Generally, you'll only have to do this once after you first clone the project.",
			Action:      bootstrap,
		},
	}

	app.Run(os.Args)
}
