package main

import (
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
)

func build(c *cli.Context) {
	consfile, err := getConsfile()
	if err != nil {
		errfAndExit(1, "error getting consfile [%s]", err)
	}
	args := []string{"build"}
	if consfile.Build.Output != "" {
		args = append(args, "-o")
		args = append(args, consfile.Build.Output)
	}
	cmd := exec.Command("go", args...)
	cmd.Env = os.Environ()
	for _, envStr := range consfile.Build.Env {
		cmd.Env = append(cmd.Env, envStr)
	}
	cmd.Dir = "" // force using the current working directory

	statusf(cmdStr(cmd))

	out, err := cmd.CombinedOutput()
	if err != nil {
		errfAndExit(1, string(out))
	}
	if len(out) == 0 {
		successf("success")
	} else {
		successf(string(out))
	}
}
