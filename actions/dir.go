package actions

import (
	"fmt"
	"path/filepath"
	"strings"
)

// packagePath gets the portion of fullPath under gopath/src
func packagePath(gopath, fullPath string) (string, error) {
	srcPath := filepath.Join(gopath, "src")
	if !strings.HasPrefix(fullPath, srcPath) {
		return "", fmt.Errorf("%s is not under the gopath 'src' dir (%s)", fullPath, srcPath)
	}
	return fullPath[len(srcPath):], nil
}
