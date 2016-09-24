package actions

import (
	"testing"

	"github.com/arschles/assert"
)

func TestPackagePath(t *testing.T) {
	gopath := "/go"
	full := "/go/src/github.com/arschles/godo"

	pkg, err := packagePath(gopath, full)
	assert.NoErr(t, err)
	assert.Equal(t, pkg, "github.com/arschles/godo", "package")
}
