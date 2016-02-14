package server

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
)

type errCreatingHeader struct {
	path string
	err  error
}

func (e errCreatingHeader) Error() string {
	return fmt.Sprintf("error creating header for %s (%s)", e.path, e.err)
}

func createTarArchive(wr io.Writer, files []string) (*tar.Writer, error) {
	tarWriter := tar.NewWriter(wr)
	for _, file := range files {
		fd, fInfo, err := openAndStatFile(file)
		if err != nil {
			return nil, err
		}
		hdr := &tar.Header{Name: file, Mode: int64(fInfo.Mode()), Size: fInfo.Size()}
		if err := tarWriter.WriteHeader(hdr); err != nil {
			return nil, errCreatingHeader{path: file, err: err}
		}

		fileBytes, err := ioutil.ReadAll(fd)
		if err != nil {
			return nil, err
		}
		if _, err := tarWriter.Write(fileBytes); err != nil {
			return nil, err
		}
	}
	if err := tarWriter.Close(); err != nil {
		return nil, err
	}
	return tarWriter, nil
}
