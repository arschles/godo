package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

const defaultConsFileName = "gocons.yaml"

func main() {
	if len(os.Args) != 2 {
		errAndExit(1, "missing build directive")
	}
	cmd := os.Args[1]

	consBytes, err := ioutil.ReadFile(defaultConsFileName)
	if err != nil {
		errAndExit(1, "error reading cons file [%s]", err)
	}

	consFile := &Consfile{}
	if err := yaml.Unmarshal(consBytes, consFile); err != nil {
		errAndExit(1, "bad consfile [%s]", err)
	}

	if cmd == "build" {
		args := []string{"build"}
		if consFile.Build.Output != "" {
			args = append(args, "-o")
			args = append(args, consFile.Build.Output)
		}
		statusf("go %s", strings.Join(args, " "))
		cmd := exec.Command("go", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			errAndExit(1, string(out))
		}
		if len(out) == 0 {
			successf("success")
		} else {
			successf(string(out))
		}
	} else {
		errAndExit(1, "unsupported command %s", cmd)
	}
}
