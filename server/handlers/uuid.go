package handlers

import (
	"github.com/pborman/uuid"
)

func shortUUID() string {
	return uuid.New()[0:7]
}
