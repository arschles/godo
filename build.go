package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/codegangsta/cli"
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
	os.Setenv("GO15VENDOREXPERIMENT", "1")
	cmd.Env = append(cmd.Env, "GO15VENDOREXPERIMENT=1")
	cmd.Dir = "" // force using the current working directory
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
