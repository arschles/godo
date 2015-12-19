package build

type Pipeline struct {
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"` // can be templated
	Steps       []PipelineStep `yaml:"steps"`
}
type PipelineStep struct {
	Name   string              `yaml:"name"` // can be templated
	Params []PipelineStepParam `yaml:"params"`
}
type PipelineStepParam struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"` // can be templated
}
