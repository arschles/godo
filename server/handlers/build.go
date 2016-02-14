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
		dir := filepath.Dir(fileName)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Err("Creating base directory %s (%s)", dir, err)
			continue
		}
		fd, err := os.Create(fileName)
		if err != nil {
			log.Err("Creating %s (%s)", fileName, err)
			continue
		}
		// TODO: fix too many open files. wrap this operation in a func so it can be deferred
		defer fd.Close()

		var dest io.Writer = fd
		isGCIFile := hdr.Name == config.DefaultFileNameYaml || hdr.Name == config.DefaultFileNameYml
		if isGCIFile && gciFileName == "" {
			dest = io.MultiWriter(fd, &gciFileBytes)
			gciFileName = hdr.Name
		} else if isGCIFile && gciFileName != "" {
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
