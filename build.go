package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/codegangsta/cli"
)

func cmdStr(cmd *exec.Cmd) string {
	var cmds []string
	for _, arg := range cmd.Args {
		cmds = append(cmds, arg)
	}

	return fmt.Sprintf("%s (Env %s", strings.Join(cmds, " "), strings.Join(cmd.Env, ":"))
}

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
	cmd := exec.Command("go", args...)
	for name, val := range consfile.Build.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", name, val))
	}
	cmd.Dir = "" // force using the current working directory

	statusf(cmdStr(cmd))

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
