package server

import (
	"archive/tar"
	"os"

	"github.com/arschles/gci/actions"
	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	"github.com/arschles/gci/server/rpc"
	"github.com/codegangsta/cli"
)

func Build(c *cli.Context) {
	cfg := config.ReadOrDie(c.String(actions.FlagConfigFile))
	cl := rpc.NewHTTPClient(cfg.CI.Build.GetHost(), cfg.CI.Build.GetPort())
	// assemble tar archive here
	var tarArchive *tar.Reader
	res, err := cl.Build(tarArchive, cfg.CI.Build.CrossCompile, cfg.CI.Build.Env)
	if err != nil {
		log.Err("building on the server (%s)", err)
		os.Exit(1)
	}
	defer res.Close()
	// read tar archive here
}
