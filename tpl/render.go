package tpl

import (
	"bytes"
	"text/template"
)

func Render(name, str string, data interface{}) (string, error) {
	t, err := template.New(name).Funcs(Funcs).Parse(str)
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	t.Execute(&out, data)
	return string(out.Bytes()), nil
}
