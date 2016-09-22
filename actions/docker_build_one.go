package actions

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	docker "github.com/fsouza/go-dockerclient"
)

type dockerBuildOneErr struct {
	err      error
	imgBuild config.ImageBuild
}

func (e dockerBuildOneErr) Error() string {
	return fmt.Sprintf("error building %s:%s (%s)", e.imgBuild.Name, e.imgBuild.GetTag(), e.err)
}

func dockerBuildOne(cl *docker.Client, imgBuild config.ImageBuild, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	dockerfileLocation := imgBuild.GetDockerfileLocation()
	dockerfileBytes, err := ioutil.ReadFile(dockerfileLocation)
	if err != nil {
		log.Err("Reading Dockerfile %s [%s]", dockerfileLocation, err)
		errCh <- dockerBuildOneErr{err: err, imgBuild: imgBuild}
		return
	}

	t := time.Now()
	buf := bytes.NewBuffer(nil)
	tr := tar.NewWriter(buf)
	tr.WriteHeader(&tar.Header{
		Name:       "Dockerfile",
		Size:       int64(len(dockerfileBytes)),
		ModTime:    t,
		AccessTime: t,
		ChangeTime: t,
	})
	tr.Write(dockerfileBytes)

	dir := imgBuild.Context.GetDirectory()
	buildCtx, err := filepath.Abs(dir)
	if err != nil {
		log.Err("Invalid Docker build context %s [%s]", dir, err)
		errCh <- dockerBuildOneErr{err: err, imgBuild: imgBuild}
		return
	}

	skipSet := imgBuild.Context.GetSkips()
	err = tarDir(buildCtx, tr, func(path string, fi os.FileInfo) bool {
		if _, ok := skipSet[path]; ok {
			return true
		}
		if strings.Contains(path, ".git") {
			return true
		}
		if fi.Name() == "Dockerfile" {
			return true
		}
		return false
	})
	if err != nil {
		log.Err("Archiving the build context directory %s [%s]", buildCtx, err)
		errCh <- dockerBuildOneErr{err: err, imgBuild: imgBuild}
	}

	if err := tr.Close(); err != nil {
		log.Err("Closing the build context archive preparing to send it to the Docker daemon [%s]", err)
		errCh <- dockerBuildOneErr{err: err, imgBuild: imgBuild}
	}

	opts := docker.BuildImageOptions{
		Name:           fmt.Sprintf("%s:%s", imgBuild.Name, imgBuild.GetTag()),
		InputStream:    buf,
		Dockerfile:     "Dockerfile",
		OutputStream:   os.Stdout,
		RmTmpContainer: true,
		Pull:           true,
	}
	if err := cl.BuildImage(opts); err != nil {
		log.Err("Building image %s [%s]", imgBuild.Name, err)
		errCh <- dockerBuildOneErr{err: err, imgBuild: imgBuild}
	}
	log.Info("Successfully built Docker image %s", opts.Name)
}
