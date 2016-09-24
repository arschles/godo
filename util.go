package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/arschles/godo/log"
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

// execOrDie prints the command before executing it, executes it and returns its output.
// if the command failed, calls log.Die with a helpful error message
func runOrDie(cmd *exec.Cmd, env []string) []byte {
	cmd.Env = env

	log.Info(cmdStr(cmd))
	log.Debug("Env: %s", envStr(cmd))

	out, err := cmd.CombinedOutput()
	if err != nil {
		var s string
		if len(out) > 0 {
			s += fmt.Sprintf("%s\n", string(out))
		}
		if err != nil {
			s += err.Error()
		}
		log.Die(s)
	}
	return out
}
