package file

import (
	"io"
	"os"
)

// CreateAndWrite creates a file called "name" (or overwrites if it already existed), reads all of r into the file, then closes the file. Returns any errors encountered along the way.
func CreateAndWrite(name string, r io.Reader) error {
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()
	if _, err := io.Copy(fd, r); err != nil {
		return err
	}
	return nil
}
