package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/arschles/gocons/log"
	"github.com/codegangsta/cli"
)

func run(c *cli.Context) {
	consfile := getConsfileOrDie()
	tgtName := c.Args().First()
	if tgtName == "" {
		log.Die("no target given")
	}
	var tgt *Target = nil
	for _, target := range consfile.Targets {
		if target.Name == tgtName {
			tgt = &target
			break
		}
	}
	if tgt == nil {
		log.Die("no target %s", tgtName)
	}
	for _, cmd := range tgt.Commands {
		if cmd == "" {
			log.Die("command %s is empty", cmd)
		}
		cmdSpl := strings.Split(cmd, " ")
		cmd := exec.Command(cmdSpl[0], cmdSpl[1:]...)
		out := runOrDie(cmd, append(os.Environ(), Envs(consfile.Envs).Strings()...))
		if len(out) > 0 {
			log.Info(string(out))
		}
	}
	log.Info("done")
}
