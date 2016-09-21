package actions

import (
	"os"

	"github.com/arschles/steward/log"
	"github.com/fsouza/go-dockerclient"
)

func getAuthConfigs(authFileLoc string) (*docker.AuthConfigurations, func(), error) {
	authFile, err := os.Open(authFileLoc)
	if err != nil {
		log.Err("Reading Docker auth file %s [%s]", authFileLoc, err)
		return nil, nil, err
	}
	retFn := func() {
		if err := authFile.Close(); err != nil {
			log.Err("Closing Docker auth file %s [%s]", authFileLoc, err)
		}
	}

	auths, err := docker.NewAuthConfigurations(authFile)
	if err != nil {
		log.Err("Parsing auth file %s [%s]", authFileLoc, err)
		return nil, nil, err
	}
	return auths, retFn, nil
}
