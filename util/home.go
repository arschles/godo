package util

import (
	"os"
)

const (
	DockerImageConst = "DOCKERIMAGE"
)

func GetHome() string {
	if os.Getenv(DockerImageConst) == "true" {
		return "/dockerhome"
	}
	return os.Getenv("HOME")
}
