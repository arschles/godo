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
	"github.com/arschles/gci/server/common"
	dockutil "github.com/arschles/gci/util/docker"
	fileutil "github.com/arschles/gci/util/file"
	tarutil "github.com/arschles/gci/util/tar"
	docker "github.com/fsouza/go-dockerclient"
)

type test struct {
	dockerCl      *docker.Client
	tmpDirCreator fileutil.TmpDirCreator
}

func NewTest(dockerCl *docker.Client, tmpDirCreator fileutil.TmpDirCreator) http.Handler {
	return &test{
		dockerCl:      dockerCl,
		tmpDirCreator: tmpDirCreator,
	}
}

func (b *test) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	packageName := r.Header.Get(common.PackageNameHeader)
	if packageName == "" {
		http.Error(w, fmt.Sprintf("You must include a %s header", common.PackageNameHeader), http.StatusBadRequest)
		return
	}
	tr := tar.NewReader(r.Body)
	defer r.Body.Close()

	buildUUID := shortUUID()
	srcTmpDir, err := b.tmpDirCreator("%s-src-%s", tmpDirPrefix, buildUUID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating temp directory for source files (%s)", err), http.StatusInternalServerError)
		return
	}
	binTmpDir, err := b.tmpDirCreator("%s-bin-%s", tmpDirPrefix, buildUUID)
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

	logsCh := make(chan string)
	resultCh := make(chan int)
	errCh := make(chan error)
	if cfg.Build.OutputBinary == "" {
		cfg.Build.OutputBinary = filepath.Base(packageName)
	}

	go dockutil.Test(b.dockerCl, srcTmpDir, packageName, containerGoPath, cfg, logsCh, resultCh, errCh)

	flush := func() {}
	if fl, ok := w.(http.Flusher); ok {
		flush = func() { fl.Flush() }
	}
	for {
		select {
		case l := <-logsCh:
			w.Write([]byte(l))
			flush()
		case err := <-errCh:
			http.Error(w, fmt.Sprintf("Error running tests (%s)", err), http.StatusInternalServerError)
			return
		case code := <-resultCh:
			if code == 0 {
				return
			} else {
				w.Write([]byte(fmt.Sprintf("Tests exited with code %d", code)))
			}
		}
	}
}
