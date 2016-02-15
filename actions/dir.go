package actions

import (
	"path/filepath"
)

// packagePath gets the portion of fullPath under gopath/src
func packagePath(gopath, fullPath string) (string, error) {
	return filepath.Rel(gopath+"/src", fullPath)
}
