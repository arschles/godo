package rpc

import (
	"io"
)

type Client interface {
	// Build sends the build context to the server, along with the cross compilation directive and the environment.
	// The return value is either a tar archive of the results or an error indicating what went wrong
	Build(buildContext io.Reader, crossCompile bool, env []string) (io.ReadCloser, error)
}
