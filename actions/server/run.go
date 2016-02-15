package server

import (
	"fmt"
	"net/http"

	"github.com/arschles/gci/actions"
	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	"github.com/arschles/gci/server/handlers"
	dockutil "github.com/arschles/gci/util/docker"
	fileutil "github.com/arschles/gci/util/file"
	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
)

func Run(c *cli.Context) {
	dockerCl := dockutil.ClientOrDie()
	mux := http.NewServeMux()
	cfg := config.ReadOrDie(c.String(actions.FlagConfigFile))
	mux.Handle("/build", handlers.NewBuild(dockerCl, fileutil.LocalTmpDirCreator())).Methods("POST")
	mux.Handle("/test", handlers.NewTest(dockerCl, fileutil.LocalTmpDirCreator())).Methods("POST")
	hostStr := fmt.Sprintf("%s:%d", cfg.CI.Server.GetBindHost(), cfg.CI.Server.GetPort())
	log.Info("Serving GCI on %s", hostStr)
	log.Die(http.ListenAndServe(hostStr, mux).Error())
}
