package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	fileutil "github.com/arschles/gci/util/file"
)

type ErrCreatingHeader struct {
	path string
	err  error
}

func (e ErrCreatingHeader) Error() string {
	return fmt.Sprintf("error creating header for %s (%s)", e.path, e.err)
}

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

func WriteFileToTarWriter(file string, tarWriter *tar.Writer) error {
	fd, fInfo, err := fileutil.OpenAndStat(file)
	if err != nil {
		return err
	}
	defer fd.Close()
	hdr := &tar.Header{Name: file, Mode: int64(fInfo.Mode()), Size: fInfo.Size()}
	if err := tarWriter.WriteHeader(hdr); err != nil {
		return ErrCreatingHeader{path: file, err: err}
	}
	fileBytes, err := ioutil.ReadAll(fd)
	if err != nil {
		return err
	}
	if _, err := tarWriter.Write(fileBytes); err != nil {
		return err
	}
	return nil
}

// CreateArchiveFromFiles creates a tar archive from the given files and writes it into wr. On any non-nil error, the contents of wr will be undefined
func CreateArchiveFromFiles(wr io.Writer, files []string) error {
	tarWriter := tar.NewWriter(wr)
	for _, file := range files {
		if err := WriteFileToTarWriter(file, tarWriter); err != nil {
			return err
		}
	}
	if err := tarWriter.Close(); err != nil {
		return err
	}
	return nil
}
