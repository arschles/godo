package server

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/arschles/gci/config/ci"
)

// an os.FileInfo implementation, just for testing
type testFileInfo struct {
	name  string
	isDir bool
}

func (t testFileInfo) Name() string       { return filepath.Base(t.name) }
func (t testFileInfo) Size() int64        { return 0 }
func (t testFileInfo) Mode() os.FileMode  { return os.ModePerm }
func (t testFileInfo) ModTime() time.Time { return time.Now() }
func (t testFileInfo) IsDir() bool        { return t.isDir }
func (t testFileInfo) Sys() interface{}   { return nil }

func TestMatchesExcludeRecursive(t *testing.T) {
	path := "/a/b/c"
	info := testFileInfo{name: "c", isDir: false}
	excludes := []ci.Exclude{ci.Exclude{Name: "c", Recursive: true}}
	assert.True(t, matchesExclude(path, info, excludes), "no match when expected")
	excludes[0].Recursive = false
	assert.False(t, matchesExclude(path, info, excludes), "match when not expected")
}

func TestMatchesExcludeNonRecursive(t *testing.T) {
	path := "/a/b/c"
	info := testFileInfo{name: "c", isDir: false}
	excludes := []ci.Exclude{ci.Exclude{Name: "/a/b/c", Recursive: false}}
	assert.True(t, matchesExclude(path, info, excludes), "no match when expected")
}
