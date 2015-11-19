package main

import (
	"os/exec"
	"strings"
)

func cmdStr(cmd *exec.Cmd) string {
	var cmds []string
	for _, arg := range cmd.Args {
		cmds = append(cmds, arg)
	}
	return strings.Join(cmds, " ")
}

func envStr(cmd *exec.Cmd) string {
	return strings.Join(cmd.Env, ":")
}
