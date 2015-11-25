package main

import (
	"os"
	"os/exec"

	"github.com/arschles/gocons/log"
	"github.com/codegangsta/cli"
)

func build(c *cli.Context) {
	consfile, err := getConsfile()
	if err != nil {
		log.Die("error getting consfile [%s]", err)
	}
	args := []string{"build"}
	if consfile.Build.Output != "" {
		args = append(args, "-o")
		args = append(args, consfile.Build.Output)
	}
	cmd := exec.Command("go", args...)
	env := os.Environ()
	cmd.Dir = "" // force using the current working directory
	for _, envStr := range consfile.Build.Env {
		env = append(env, envStr)
	}
	out := runOrDie(cmd, env)
	if len(out) == 0 {
		log.Info("done")
	} else {
		log.Info(string(out))
	}
}
