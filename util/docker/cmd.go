package docker

import (
	"fmt"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

// CmdStr returns the 'docker run' command that you'd execute to achieve the run configuration
// represented by co and hc
func CmdStr(co docker.CreateContainerOptions, hc docker.HostConfig) string {
	ret := []string{"docker run"}
	for _, env := range co.Config.Env {
		ret = append(ret, fmt.Sprintf("-e %s", env))
	}

	for _, b := range hc.Binds {
		ret = append(ret, fmt.Sprintf("-v %s", b))
	}
	ret = append(ret, fmt.Sprintf("-w %s", co.Config.WorkingDir))
	ret = append(ret, fmt.Sprintf("--name=%s", co.Name))
	ret = append(ret, co.Config.Image)
	ret = append(ret, strings.Join(co.Config.Cmd, " "))

	return strings.Join(ret, " ")
}
