package main

import (
	"fmt"
	"strings"

	"github.com/arschles/gocons/log"
	"github.com/codegangsta/cli"
)

func targets(c *cli.Context) {
	consfile, err := getConsfile()
	if err != nil {
		log.Die("error getting consfile [%s]", err)
	}

	var s []string
	if consfile.Bootstrap != nil {
		s = append(s, fmt.Sprintf("bootstrap - bootstrap this project"))
	}
	if consfile.Build != nil {
		s = append(s, fmt.Sprintf("build - build this project"))
	}
	if consfile.Test != nil {
		s = append(s, fmt.Sprintf("test - test this project"))
	}
	if consfile.Install != nil {
		s = append(s, fmt.Sprintf("install - install this project"))
	}
	for _, otherCmd := range consfile.OtherCommands {
		s = append(s, fmt.Sprintf("%s - %s", otherCmd.Name, otherCmd.Description))
	}

	log.Msg(strings.Join(s, "\n"))
}
