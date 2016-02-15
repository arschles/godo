package tar

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/arschles/assert"
)

func TestNewFile(t *testing.T) {
	path, name := "/a/b/c", "myfile"
	f := NewFile(path, name)
	assert.Equal(t, f.Name(), name, "name")
	assert.Equal(t, f.Path(), path, "path")
}

func TestFilesFromRoot(t *testing.T) {
	root := "/myroot"
	baseNames := []string{"a", "b/c/d"}
	files := FilesFromRoot(root, baseNames, filepath.Join)
	assert.Equal(t, len(files), len(baseNames), "number of files")
	for _, file := range files {
		if !strings.HasPrefix(file.Path(), root) {
			t.Errorf("file %s (%s) didn't have expected prefix %s", file.Name(), file.Path(), root)
		}
	}
}
