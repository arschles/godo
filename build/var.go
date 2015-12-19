package build

type Var struct {
	Name    string `yaml:"name"`
	Default string `yaml:"default"`
	Env     string `yaml:"env"`
}
