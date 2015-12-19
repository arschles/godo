package build

type Pipeline struct {
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"`
	Steps       []PipelineStep `yaml:"steps"`
}
type PipelineStep struct {
	Name   string              `yaml:"name"`
	Params []PipelineStepParam `yaml:"params"`
}
type PipelineStepParam struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
