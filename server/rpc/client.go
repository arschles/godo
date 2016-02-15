package rpc

import (
	"io"
)

type Client interface {
	// Build sends the build context to the server, along with the cross compilation directive and the environment.
	// On success, writes the results to target. On failure, returns an indicative error.
	// Note that writes may have happened on target even in the error case.
	Build(buildContext io.Reader, target io.Writer, crossCompile bool, packageName string, env []string) error

	// Test sends the build context to the server along with environment variables.
	// On success, sends logs to target. On failure, returns an indicative error.
	// Writes may have happened on target even in the error case
	Test(buildContext io.Reader, target io.Writer, packageName string, env []string) error
}
