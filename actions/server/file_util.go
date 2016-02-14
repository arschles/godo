package server

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/arschles/gci/config/ci"
)

type errOpenFile struct {
	path string
	err  error
}

func (e errOpenFile) Error() string {
	return fmt.Sprintf("Error opening %s (%s)", e.path, e.err)
}

type errStatFile struct {
	path string
	err  error
}

func (e errStatFile) Error() string {
	return fmt.Sprintf("Error getting info for %s (%s)", e.path, e.err)
}

// openandStatFile calls os.Opena and os.Stat on path, returning the result of both calls if both were successful.
//
// If either returned a non-nil error, returns errOpenFile or errStatFile (respectively) and nil for both the File and the FileInfo
func openAndStatFile(path string) (*os.File, os.FileInfo, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, nil, errOpenFile{path: path, err: err}
	}
	fInfo, err := os.Stat(path)
	if err != nil {
		return nil, nil, errStatFile{path: path, err: err}
	}
	return fd, fInfo, nil
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

// getFiles walks the directory staring at root. it returns all of the files (as relative paths) it finds, skipping all files and directories in exceptions
func getFiles(root string, excludes []ci.Exclude) ([]string, error) {
	var paths []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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

// writeToFile creates a file called name and writes the contents of rc to it.
// If name already existed, it erases it first
func writeToFile(name string, rc io.Reader) error {
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()
	if _, err := io.Copy(fd, rc); err != nil {
		return err
	}
	return nil
}
