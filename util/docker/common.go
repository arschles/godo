package docker

import (
	"fmt"
)

func ContainerGopath(packageName string) string {
	return fmt.Sprintf("/go/src/%s", packageName)
}
