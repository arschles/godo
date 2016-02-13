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
	gciFileName := ""
	var gciFileBytes bytes.Buffer
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
		defer fd.Close()

		var dest io.Writer = fd
		if hdr.Name == config.DefaultFileNameYaml || hdr.Name == config.DefaultFileNameYml && gciFileName != "" {
			dest = io.MultiWriter(fd, &gciFileBytes)
			gciFileName = hdr.Name
		} else {
			str := fmt.Sprintf("multiple GCI config files found (%s and %s)", hdr.Name, gciFileName)
			http.Error(w, str, http.StatusBadRequest)
			return
		}

		if _, err := io.Copy(dest, tr); err != nil {
			log.Err("Writing archived file to %s (%s)", fileName, err)
			continue
		}
	}
	cfg := config.Empty()
	if gciFileName != "" {
		c, err := config.ReadBytes(gciFileBytes.Bytes())
		if err != nil {
			http.Error(w, fmt.Sprintf("%s was invalid (%s)", gciFileName, err), http.StatusBadRequest)
			return
		}
		cfg = c
	}
	io.Copy(w, strings.NewReader(fmt.Sprintf("%s", *cfg)))
}
