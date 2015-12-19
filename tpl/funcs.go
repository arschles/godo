package tpl

import (
	"os"
	"text/template"
)

var Funcs = template.FuncMap(map[string]interface{}{
	"PWD": func() (string, error) { return os.Getwd() },
})
