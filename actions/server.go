package actions

import (
	"fmt"
	"net/http"

	"github.com/arschles/gci/config"
	"github.com/arschles/gci/log"
	"github.com/arschles/gci/server/handlers"
	"github.com/codegangsta/cli"
)

func Server(c *cli.Context) {
	mux := http.NewServeMux()
	cfg := config.ReadOrDie(c.String(FlagConfigFile))
	mux.Handle("/build", handlers.Build{})
	hostStr := fmt.Sprintf("%s:%d", cfg.CI.Server.GetHost(), cfg.CI.Server.GetPort())
	log.Info("Serving GCI on %s", hostStr)
	log.Die(http.ListenAndServe(hostStr, mux).Error())
}
