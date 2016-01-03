package actions

import (
	"os"

	"github.com/arschles/gci/log"
)

type paths struct {
	gopath string
	cwd    string
	pkg    string
}

func pathsOrDie() paths {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Err("GOPATH environment variable not found")
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Err("getting current working dir (%s)", err)
		os.Exit(1)
	}

	pkgPath, err := packagePath(gopath, cwd)
	if err != nil {
		log.Err("Error detecting package name [%s]", err)
		os.Exit(1)
	}

	return paths{gopath: gopath, pkg: pkgPath, cwd: cwd}
}
