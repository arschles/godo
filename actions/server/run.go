package server

import (
	"fmt"
	"net/http"

	"github.com/arschles/gci/actions"
	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	"github.com/arschles/gci/server/handlers"
	dockutil "github.com/arschles/gci/util/docker"
	"github.com/codegangsta/cli"
)

func Run(c *cli.Context) {
	paths := actions.PathsOrDie()
	dockerCl := dockutil.ClientOrDie()
	mux := http.NewServeMux()
	cfg := config.ReadOrDie(c.String(actions.FlagConfigFile))
	mux.Handle("/build", handlers.NewBuild(paths.GoPath, paths.CWD, dockerCl))
	hostStr := fmt.Sprintf("%s:%d", cfg.CI.Server.GetBindHost(), cfg.CI.Server.GetPort())
	log.Info("Serving GCI on %s", hostStr)
	log.Die(http.ListenAndServe(hostStr, mux).Error())
}
