package build

import "os"

type Var struct {
	Name    string `yaml:"name"`
	Default string `yaml:"default"`
	Env     string `yaml:"env"`
}

func (v Var) GetValue() string {
	e := os.Getenv(v.Env)
	if e != "" {
		return e
	}
	return v.Default
}

type VarMap map[string]Var
