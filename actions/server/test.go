package server

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/arschles/gci/actions"
	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	"github.com/arschles/gci/server/rpc"
	fileutil "github.com/arschles/gci/util/file"
	tarutil "github.com/arschles/gci/util/tar"
	"github.com/codegangsta/cli"
	humanize "github.com/dustin/go-humanize"
)

func Test(c *cli.Context) {
	cfg := config.ReadOrDie(c.String(actions.FlagConfigFile))
	wd, err := os.Getwd()
	if err != nil {
		log.Err("getting current working directory (%s)", err)
		os.Exit(1)
	}

	paths := actions.PathsOrDie()

	log.Info("Walking current directory")

	fileBaseNames, err := fileutil.WalkAndExclude(paths.CWD, true, cfg.CI.Build.Excludes)
	if err != nil {
		log.Err("walking %s to get files to upload to the server (%s)", wd, err)
		os.Exit(1)
	}

	files := tarutil.FilesFromRoot(paths.CWD, fileBaseNames, filepath.Join)

	log.Info("Archiving %d files", len(files))

	tarArchive := new(bytes.Buffer)
	if err := tarutil.CreateArchiveFromFiles(tarArchive, files); err != nil {
		log.Err("creating tar archive of %s (%s)", wd, err)
		os.Exit(1)
	}

	log.Info("Sending %s tar archive to server", humanize.Bytes(uint64(tarArchive.Len())))
	cl := rpc.NewHTTPClient(cfg.CI.Build.GetHost(), cfg.CI.Build.GetPort())
	if err := cl.Test(tarArchive, os.Stdout, paths.PackageName, cfg.CI.Build.Env); err != nil {
		log.Err("testing on the server (%s)", err)
		os.Exit(1)
	}

	log.Info("Success")
}
