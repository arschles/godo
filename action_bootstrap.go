package main

import (
	"os"
	"os/exec"

	"github.com/arschles/gocons/log"
	"github.com/codegangsta/cli"
)

func bootstrap(c *cli.Context) {
	consfile, err := getConsfile()
	if err != nil {
		log.Die("error getting consfile [%s]", err)
	}
	for _, str := range consfile.Bootstrap.Commands {
		cmd := exec.Command(str)
		out := runOrDie(cmd, os.Environ())
		log.Info(string(out))
	}
}
