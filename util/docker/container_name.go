package docker

import (
	"fmt"
	"path/filepath"

	"github.com/pborman/uuid"
)

// NewContainerName returns a new, unique container name that includes prefix and cwd
func NewContainerName(prefix string, cwd string) string {
	projName := filepath.Base(cwd)
	return fmt.Sprintf("godo-%s-%s-%s", projName, prefix, uuid.New())
}
