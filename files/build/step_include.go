package build

type StepInclude struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type StepIncludesMap map[string]string
