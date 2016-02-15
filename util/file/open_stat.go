package file

import (
	"fmt"
	"os"
)

type ErrOpenFile struct {
	path string
	err  error
}

func (e ErrOpenFile) Error() string {
	return fmt.Sprintf("Error opening %s (%s)", e.path, e.err)
}

type ErrStatFile struct {
	path string
	err  error
}

func (e ErrStatFile) Error() string {
	return fmt.Sprintf("Error getting info for %s (%s)", e.path, e.err)
}

// openandStatFile calls os.Opena and os.Stat on path, returning the result of both calls if both were successful.
//
// If either returned a non-nil error, returns errOpenFile or errStatFile (respectively) and nil for both the File and the FileInfo
func OpenAndStat(path string) (*os.File, os.FileInfo, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, nil, ErrOpenFile{path: path, err: err}
	}
	fInfo, err := os.Stat(path)
	if err != nil {
		return nil, nil, ErrStatFile{path: path, err: err}
	}
	return fd, fInfo, nil
}
