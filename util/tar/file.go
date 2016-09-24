package tar

import (
	"archive/tar"
	"fmt"
	"io"

	fileutil "github.com/arschles/godo/util/file"
)

type ErrWritingHeader struct {
	path string
	err  error
}

func (e ErrWritingHeader) Error() string {
	return fmt.Sprintf("error creating header for %s (%s)", e.path, e.err)
}

// File represents an on-disk file that can be written by a *tar.Writer, possibly as a different name, to a tar archive
type File struct {
	loc     string
	tarName string
}

// NewFile creates a new File struct that points to the file at path on disk. It will be written as name to a tar archive on subsequent successful calls to Write
func NewFile(path, name string) *File {
	return &File{loc: path, tarName: name}
}

func FilesFromRoot(root string, baseNames []string, joiner func(elts ...string) string) []*File {
	ret := make([]*File, len(baseNames))
	for i, baseName := range baseNames {
		ret[i] = NewFile(joiner(root, baseName), baseName)
	}
	return ret
}

func (d *File) String() string {
	return fmt.Sprintf("%s (%s)", d.tarName, d.loc)
}

// Name returns the name of the file as it will appear in a tar archive after a successful call to Write
func (d *File) Name() string {
	return d.tarName
}

// Path returns the path on disk of the file
func (d *File) Path() string {
	return d.loc
}

// Write opens the file, writes it to
func (d *File) Write(tw *tar.Writer) error {
	fd, fInfo, err := fileutil.OpenAndStat(d.loc)
	if err != nil {
		return err
	}
	defer fd.Close()
	hdr := &tar.Header{Name: d.tarName, Mode: int64(fInfo.Mode()), Size: fInfo.Size()}
	if err := tw.WriteHeader(hdr); err != nil {
		return ErrWritingHeader{path: d.loc, err: err}
	}
	if _, err := io.Copy(tw, fd); err != nil {
		return err
	}
	return nil
}
