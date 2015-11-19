package main

import (
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
	"github.com/helm/helm/log"
)

func other(c *cli.Context) {
	consfile, err := getConsfile()
	if err != nil {
		log.Die("error getting consfile [%s]", err)
	}
	args := c.Args()
	if len(args) != 1 {
		log.Die("other command not specified")
	}
	otherName := args[0]
	var otherCmd OtherCommand
	found := false
	for _, other := range consfile.OtherCommands {
		if other.Name == otherName {
			otherCmd = other
			found = true
			break
		}
	}
	if !found {
		log.Die("no %s command listed in the consfile", otherName)
	}
	if otherCmd.Command == "" {
		log.Die("no command specified for %s", otherName)
	}
	cmd := exec.Command(otherCmd.Command)
	cmd.Env = os.Environ()
	log.Info(cmdStr(cmd))
	log.Debug("Env: %s", envStr(cmd))
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Die("error running command %s (%s)", otherCmd.Command, err)
	}
	log.Info(string(out))

}
