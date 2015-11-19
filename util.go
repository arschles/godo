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
	//TODO: print out the env in debug mode
}
