package tar

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func CopyToFile(tarReader *tar.Reader, fileName string, otherWriters ...io.Writer) error {
	dir := filepath.Dir(fileName)
	// TODO: better mode for this dir?
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	fd, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer fd.Close()
	var dest io.Writer = fd
	if len(otherWriters) > 0 {
		for _, otherWriter := range otherWriters {
			dest = io.MultiWriter(dest, otherWriter)
		}
	}
	if _, err := io.Copy(dest, tarReader); err != nil {
		return err
	}
	return nil
}

// CreateArchiveFromFiles creates a tar archive from the given files and writes it into wr. On any non-nil error, the contents of wr will be undefined
func CreateArchiveFromFiles(wr io.Writer, files []*File) error {
	tarWriter := tar.NewWriter(wr)
	for _, file := range files {
		if err := file.Write(tarWriter); err != nil {
			return err
		}
	}
	if err := tarWriter.Close(); err != nil {
		return err
	}
	return nil
}
