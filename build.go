package main

import (
	"github.com/codegangsta/cli"
	"os/exec"
	"strings"
)

func build(c *cli.Context) {
	consfile, err := getConsfile()
	if err != nil {
		errAndExit(1, "error getting consfile [%s]", err)
	}
	args := []string{"build"}
	if consfile.Build.Output != "" {
		args = append(args, "-o")
		args = append(args, consfile.Build.Output)
	}
	statusf("go %s", strings.Join(args, " "))
	cmd := exec.Command("go", args...)
	cmd.Env = append(cmd.Env, "GOVENDOREXPERIMENT=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		errAndExit(1, string(out))
	}
	if len(out) == 0 {
		successf("success")
	} else {
		successf(string(out))
	}
}
