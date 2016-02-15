package handlers

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/arschles/gci/config"
	dockutil "github.com/arschles/gci/util/docker"
	fileutil "github.com/arschles/gci/util/file"
	tarutil "github.com/arschles/gci/util/tar"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pborman/uuid"
)

const (
	tmpDirPrefix = "gci"
)

type build struct {
	buildDir      string
	dockerCl      *docker.Client
	tmpDirCreator fileutil.TmpDirCreator
}

func NewBuild(baseBuildDir string, dockerCl *docker.Client) http.Handler {
	return &build{buildDir: baseBuildDir, dockerCl: dockerCl, tmpDirCreator: fileutil.DefaultTmpDirCreator()}
}

func (b build) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "you must POST to this endpoint", http.StatusBadRequest)
		return
	}
	tr := tar.NewReader(r.Body)
	defer r.Body.Close()

	buildUUID := uuid.New()
	srcTmpDir, err := b.tmpDirCreator("", "%s-src-%s", tmpDirPrefix, buildUUID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating temp directory for source files (%s)", err), http.StatusInternalServerError)
		return
	}
	binTmpDir, err := b.tmpDirCreator("", "%s-bin-%s", tmpDirPrefix, buildUUID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating temp directory for binary files (%s)", err), http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := os.RemoveAll(srcTmpDir); err != nil {
			log.Printf("Error removing source temp dir %s (%s)", srcTmpDir, err)
		}
		if err := os.RemoveAll(binTmpDir); err != nil {
			log.Printf("Error removing binary temp dir %s (%s)", binTmpDir, err)
		}
	}()

	cfg := config.Empty()
	configFileFound := false
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error reading file (%s)", err)
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
		fileName := filepath.Join(srcTmpDir, hdr.Name)
		if err := tarutil.CopyToFile(tr, fileName, otherWriters...); err != nil {
			str := fmt.Sprintf("Error copying %s to a file (%s)", hdr.Name, err)
			log.Printf(str)
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
	if _, err := dockutil.Build(b.dockerCl, srcTmpDir, binTmpDir, cfg); err != nil {
		http.Error(w, fmt.Sprintf("Error building (%s)", err), http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening completed binary (%s)", err), http.StatusInternalServerError)
		return
	}
	files, err := fileutil.WalkAndExclude(binTmpDir, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing all output binaries (%s)", err), http.StatusInternalServerError)
		return
	}
	if err := tarutil.CreateArchiveFromFiles(w, files); err != nil {
		http.Error(w, fmt.Sprintf("Error creating tar archive from files (%s)", err), http.StatusInternalServerError)
		return
	}
}
