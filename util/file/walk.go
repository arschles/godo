package file

import (
	"log"
	"os"
	"path/filepath"

	"github.com/arschles/gci/config/ci"
)

// WalkAndExclude walks the directory staring at root and returns all of the files (as relative paths) it finds, skipping all files and directories in exceptions
func WalkAndExclude(root string, excludes []ci.Exclude) ([]string, error) {
	var paths []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error found in WalkAndExclude (%s)", err)
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}
		if matchesExclude(rel, info, excludes) {
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}
		if info.IsDir() {
			return nil
		}
		paths = append(paths, path)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return paths, nil
}

func matchesExclude(path string, info os.FileInfo, excludes []ci.Exclude) bool {
	// TODO: make this more efficient
	for _, exclude := range excludes {
		if exclude.Recursive && exclude.Name == info.Name() {
			return true
		} else if exclude.Name == path {
			if exclude.Name == path {
				return true
			}
		}
	}
	return false
}
