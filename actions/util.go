package actions

import (
	"os"

	"github.com/arschles/godo/log"
)

type Paths struct {
	GoPath      string
	CWD         string
	PackageName string
}

// PathsOrDie gets the current GOPATH, current working directory and package name and returns them all in the Paths struct. If any one of those values can't be obtained, this func logs the error and exits the process
func PathsOrDie() Paths {
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

	return Paths{GoPath: gopath, PackageName: pkgPath, CWD: cwd}
}
