package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/arschles/gocons/log"
	"github.com/codegangsta/cli"
)

var errEmptyCommandString = errors.New("empty command string")

func parseCmd(cmdStr string) (*exec.Cmd, error) {
	spl := strings.Split(cmdStr, " ")
	if len(spl) == 0 {
		return nil, errEmptyCommandString
	}
	return exec.Command(spl[0], spl[1:]...), nil
}

func install(c *cli.Context) {
	consfile, err := getConsfile()
	if err != nil {
		log.Die("error getting consfile [%s]", err)
	}
	inst := consfile.Install
	if inst.PreScript != "" {
		cmd, err := parseCmd(inst.PreScript)
		if err != nil {
			log.Die("pre script (%s)", err)
		}
		out := runOrDie(cmd, os.Environ())
		log.Info(string(out))
	}
	cmd := exec.Command("go", "install")
	out := runOrDie(cmd, os.Environ())
	log.Info(string(out))

	if inst.PostScript != "" {
		cmd, err := parseCmd(inst.PostScript)
		if err != nil {
			log.Die("post script (%s)", err)
		}
		out := runOrDie(cmd, os.Environ())
		log.Info(string(out))
	}
	log.Info("done")
}
