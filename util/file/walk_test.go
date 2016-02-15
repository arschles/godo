package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arschles/assert"
)

func testDataDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Abs(fmt.Sprintf("%s/../../testdata", wd))
}

func TestWalkAndExcludeStripPrefix(t *testing.T) {
	root, err := testDataDir()
	assert.NoErr(t, err)
	files, err := WalkAndExclude(root, true, nil)
	assert.NoErr(t, err)
	for _, file := range files {
		if strings.HasPrefix(file, root) {
			t.Errorf("%s didn't have its prefix (%s) removed", file, root)
		}
	}
}

func TestWalkAndExclude(t *testing.T) {
	root, err := testDataDir()
	assert.NoErr(t, err)
	files, err := WalkAndExclude(root, false, nil)
	assert.NoErr(t, err)
	for _, file := range files {
		if !strings.HasPrefix(file, root) {
			t.Errorf("%s didn't have its expected prefix (%s)", file, root)
		}
	}
}
