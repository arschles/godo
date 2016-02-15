package handlers

import (
	"testing"

	"github.com/arschles/assert"
)

func TestIsGCIFileName(t *testing.T) {
	assert.False(t, isGCIFileName("abc"), "abc was reported as a valid config file")
	assert.True(t, isGCIFileName("gci.yml"), "gci.yml wasn't reported as a valid config file")
	assert.True(t, isGCIFileName("gci.yaml"), "gci.yaml wasn't reported as a valid config file")
}
