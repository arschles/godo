package handlers

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/arschles/gci/log"
)

type Build struct{}

func (b Build) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "you must POST to this endpoint", http.StatusBadRequest)
		return
	}
	tr := tar.NewReader(r.Body)
	defer r.Body.Close()
	tmpDir, err := ioutil.TempDir("", "gci_server_builds")
	if err != nil {
		http.Error(w, "error creating temp directory", http.StatusBadRequest)
		return
	}
	defer os.Remove(tmpDir)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Err("Reading file (%s)", err)
			break
		}
		fileName := filepath.Join(tmpDir, hdr.Name)
		fd, err := os.Create(fileName)
		if err != nil {
			log.Err("Creating %s (%s)", fileName, err)
			continue
		}
		if _, err := io.Copy(fd, tr); err != nil {
			log.Err("Writing archived file to %s (%s)", fileName, err)
			continue
		}
	}
}
