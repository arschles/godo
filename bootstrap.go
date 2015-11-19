package main

import (
	"fmt"
	"os/exec"

	"github.com/codegangsta/cli"
)

func bootstrap(c *cli.Context) {
	consfile, err := getConsfile()
	if err != nil {
		errfAndExit(1, "error getting consfile [%s]", err)
	}
	for i, cmdStr := range consfile.Bootstrap.Commands {
		cmd := exec.Command(cmdStr)
		statusf(cmdStr)
		out, err := cmd.CombinedOutput()
		if err != nil {
			s := fmt.Sprintf("error running command %d, stopping [%s]", i+1, err)
			errfAndExit(1, s)
		}
		successf(string(out))
	}
}
