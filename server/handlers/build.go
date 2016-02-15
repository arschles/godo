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
	"time"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/server/common"
	dockutil "github.com/arschles/gci/util/docker"
	dockbuild "github.com/arschles/gci/util/docker/build"
	fileutil "github.com/arschles/gci/util/file"
	tarutil "github.com/arschles/gci/util/tar"
	docker "github.com/fsouza/go-dockerclient"
)

const (
	tmpDirPrefix = "gci"
)

type build struct {
	dockerCl      *docker.Client
	tmpDirCreator fileutil.TmpDirCreator
}

func NewBuild(dockerCl *docker.Client, tmpDirCreator fileutil.TmpDirCreator) http.Handler {
	return &build{
		dockerCl:      dockerCl,
		tmpDirCreator: tmpDirCreator,
	}
}

func (b *build) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	logsCh := make(chan dockbuild.Log)
	resultCh := make(chan int)
	errCh := make(chan error)
	if cfg.Build.OutputBinary == "" {
		cfg.Build.OutputBinary = filepath.Base(packageName)
	}

	go dockutil.Build(b.dockerCl, srcTmpDir, binTmpDir, packageName, containerGoPath, cfg, logsCh, resultCh, errCh)

	for {
		select {
		//TODO: stream logs to client!
		case l := <-logsCh:
			log.Println(l.Message())
		case err := <-errCh:
			http.Error(w, fmt.Sprintf("Error building (%s)", err), http.StatusInternalServerError)
			return
		case code := <-resultCh:
			if code == 0 {
				fileBaseNames, err := fileutil.WalkAndExclude(binTmpDir, true, nil)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error listing all output binaries (%s)", err), http.StatusInternalServerError)
					return
				}
				files := tarutil.FilesFromRoot(binTmpDir, fileBaseNames, filepath.Join)
				fmt.Println("creating archive for", files)
				tarFileName := binTmpDir + "/result.tar"
				tarFile, err := os.Create(tarFileName)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error creating intermediate tar file (%s)", err), http.StatusInternalServerError)
					return
				}
				defer tarFile.Close()
				if err := tarutil.CreateArchiveFromFiles(tarFile, files); err != nil {
					http.Error(w, fmt.Sprintf("Error creating tar archive from files (%s)", err), http.StatusInternalServerError)
					return
				}
				http.ServeContent(w, r, tarFileName, time.Now(), tarFile)
				return
			} else {
				http.Error(w, fmt.Sprintf("Build failed with code %d", code), http.StatusInternalServerError)
				return
			}
		}
	}
}
