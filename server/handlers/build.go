package handlers

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
)

const (
	tempDirPrefix = "gci_server_build"
)

type build struct {
	buildDir string
}

func NewBuild(baseBuildDir string) http.Handler {
	return &build{buildDir: baseBuildDir}
}

func (b build) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "you must POST to this endpoint", http.StatusBadRequest)
		return
	}
	tr := tar.NewReader(r.Body)
	defer r.Body.Close()
	tmpDir, err := ioutil.TempDir(b.buildDir, tempDirPrefix)
	if err != nil {
		log.Err("creating temp directory under %s (%s)", b.buildDir, err)
		http.Error(w, "error creating temp directory", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Err("Removing temporary build dir %s (%s)", tmpDir, err)
		}
	}()

	cfg := config.Empty()
	configFileFound := false
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Err("Reading file (%s)", err)
			break
		}
		var gciFileBytes bytes.Buffer
		var otherWriters []io.Writer
		if isGCIFileName(hdr.Name) && !configFileFound {
			otherWriters = append(otherWriters, &gciFileBytes)
		} else if isGCIFileName(hdr.Name) && configFileFound {
			http.Error(w, "Multiple GCI config files found", http.StatusBadRequest)
			return
		}
		fileName := filepath.Join(tmpDir, hdr.Name)
		if err := copyToFile(tr, fileName, otherWriters...); err != nil {
			str := fmt.Sprintf("Copying %s to a file (%s)", hdr.Name, err)
			log.Err(str)
			http.Error(w, str, http.StatusInternalServerError)
			return
		}
		c, err := config.ReadBytes(gciFileBytes.Bytes())
		if err != nil {
			http.Error(w, fmt.Sprintf("%s was an invalid config file (%s)", hdr.Name, err), http.StatusBadRequest)
			return
		}
		cfg = c
	}
	io.Copy(w, strings.NewReader(fmt.Sprintf("%s", *cfg)))
}
