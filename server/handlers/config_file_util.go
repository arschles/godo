package handlers

import (
	"github.com/arschles/gci/config"
)

func isGCIFileName(name string) bool {
	return name == config.DefaultFileNameYaml || name == config.DefaultFileNameYml
}
