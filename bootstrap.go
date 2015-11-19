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
	for i, str := range consfile.Bootstrap.Commands {
		cmd := exec.Command(str)
		cmd.Env = os.Environ()

		log.Info(cmdStr(cmd))
		log.Debug("Env: %s", envStr(cmd))

		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Die("error running command %d, stopping (%s)", i+1, err)
		}
		log.Info(string(out))
	}
}
