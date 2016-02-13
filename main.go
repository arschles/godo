package main

import (
	"os"

	"github.com/arschles/gci/actions"
	"github.com/arschles/gci/actions/server"
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
			Name:  actions.FlagConfigFile,
			Usage: "Specify the name and location of the config file",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "build",
			Aliases:     []string{"b"},
			Usage:       "Build your project",
			Description: "This command runs the equivalent of 'go build -o $CURRENT_DIR_NAME'",
			Action:      actions.Build,
		},
		{
			Name:        "test",
			Aliases:     []string{"t"},
			Usage:       "Test your project",
			Description: "This command runs the equivalent of 'go test ./...'",
			Action:      actions.Test,
		},
		{
			Name:        "docker-build",
			Aliases:     []string{"db"},
			Usage:       "Build a Docker image for your project",
			Description: "This command runs the equivalent of 'docker build -t $IMG_NAME $DOCKERFILE_DIR'",
			Action:      actions.DockerBuild,
		},
		{
			Name:        "docker-push",
			Aliases:     []string{"dp"},
			Usage:       "Push the Docker image for your project",
			Description: "This command runs the equivalent of 'docker push $IMG_NAME'",
			Action:      actions.DockerPush,
		},
		{
			Name:        "server",
			Aliases:     []string{"srv"},
			Usage:       "Run or access the GCI server",
			Description: "This command has subcommands to run or talk to the GCI server, all according to config parameters",
			Subcommands: []cli.Command{
				{
					Name:    "run",
					Aliases: []string{"r"},
					Usage:   "Run the GCI server according to the config file",
					Action:  server.Run,
				},
				{
					Name:        "build",
					Aliases:     []string{"b"},
					Usage:       "Build this project on a running GCI server",
					Description: "Send this project to a running GCI server and tell it to build according to the config file under ci -> build",
					Action:      server.Build,
				},
				{
					Name:        "test",
					Aliases:     []string{"t"},
					Usage:       "Test this project on a running GCI server",
					Description: "Send this project to a running GCI server and tell it to run tests according to the config file under ci -> test",
					Action:      server.Test,
				},
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		log.IsDebugging = c.Bool(actions.FlagDebug)
		return nil
	}

	app.Run(os.Args)
}
