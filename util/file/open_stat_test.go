package file

import (
	"os"
	"path/filepath"
	"testing"
	"time"
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

func TestOpenAndStat(t *testing.T) {
	t.Skip("TODO")
}
