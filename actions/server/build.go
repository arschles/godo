package server

import (
	"bytes"
	"os"

	"github.com/arschles/gci/actions"
	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	"github.com/arschles/gci/server/rpc"
	fileutil "github.com/arschles/gci/util/file"
	tarutil "github.com/arschles/gci/util/tar"
	"github.com/codegangsta/cli"
)

const (
	defaultTarOutputFile = "gci-server-build.tar"
)

func Build(c *cli.Context) {
	tarOutputFile := c.String("output")
	if tarOutputFile == "" {
		tarOutputFile = defaultTarOutputFile
	}

	cfg := config.ReadOrDie(c.String(actions.FlagConfigFile))
	wd, err := os.Getwd()
	if err != nil {
		log.Err("getting current working directory (%s)", err)
		os.Exit(1)
	}

	paths := actions.PathsOrDie()

	log.Info("Creating tar archive of current directory")

	files, err := fileutil.WalkAndExclude(wd, cfg.CI.Build.Excludes)
	if err != nil {
		log.Err("walking %s to get files to upload to the server (%s)", wd, err)
		os.Exit(1)
	}

	tarArchive := new(bytes.Buffer)
	if err := tarutil.CreateArchiveFromFiles(tarArchive, files); err != nil {
		log.Err("creating tar archive of %s (%s)", wd, err)
		os.Exit(1)
	}

	log.Info("Sending to server")
	cl := rpc.NewHTTPClient(cfg.CI.Build.GetHost(), cfg.CI.Build.GetPort())
	res, err := cl.Build(tarArchive, cfg.CI.Build.CrossCompile, paths.PackageName, cfg.CI.Build.Env)
	if err != nil {
		log.Err("building on the server (%s)", err)
		os.Exit(1)
	}
	defer res.Close()
	log.Info("Writing result to %s", tarOutputFile)
	if err := fileutil.CreateAndWrite(tarOutputFile, res); err != nil {
		log.Err("writing the result to %s (%s)", tarOutputFile, err)
		os.Exit(1)
	}
}
