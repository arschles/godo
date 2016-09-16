package docker

import (
	"fmt"
)

func ContainerGopath(gopath, packageName string) string {
	return fmt.Sprintf("%s/src/%s", gopath, packageName)
}
