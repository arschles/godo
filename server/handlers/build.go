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

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	dockutil "github.com/arschles/gci/util/docker"
	tarutil "github.com/arschles/gci/util/tar"
	docker "github.com/fsouza/go-dockerclient"
)

const (
	tempDirPrefix = "gci_server_build"
)

type build struct {
	buildDir string
	dockerCl *docker.Client
}

func NewBuild(baseBuildDir string, dockerCl *docker.Client) http.Handler {
	return &build{buildDir: baseBuildDir, dockerCl: dockerCl}
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
		if err := tarutil.CopyToFile(tr, fileName, otherWriters...); err != nil {
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
	outDir := filepath.Join(tmpDir, "built-binaries")
	if _, err := dockutil.Build(b.dockerCl, tmpDir, outDir, cfg); err != nil {
		http.Error(w, fmt.Sprintf("Error building (%s)", err), http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening completed binary (%s)", err), http.StatusInternalServerError)
		return
	}
	// TODO: tar up the outDir and send it to w
	// if _, err := io.Copy(w, []byte("hello world!")); err != nil {
	// http.Error(w, fmt.Sprintf("sending completed binary (%s)", err), http.StatusInternalServerError)
	// return
	// }
	w.Write([]byte("hello world!"))
}
