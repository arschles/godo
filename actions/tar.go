package actions

import (
	"archive/tar"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/arschles/godo/log"
)

// tarDir tars the entire directory under dir. It skips all hidden directories (starting with a '.') as well as files or directories which skip returns true for. If skip returns true for a directory, its entire contents will be skipped and skip will not be called for any of that directory's contents
//
// - All hidden directories
// - All files with names in skipNames and all directories with names in skipDirs
func tarDir(dir string, tarWriter *tar.Writer, skip func(path string, info os.FileInfo) bool) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if skip(path, info) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		log.Debug("Adding to Docker build context: %s", path)
		fileBytes, err := ioutil.ReadFile(rel)
		if err != nil {
			return err
		}
		// empty link because filepath.Walk doesn't follow symbolic links
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		tarWriter.WriteHeader(hdr)
		tarWriter.Write(fileBytes)
		return nil
	})
}
