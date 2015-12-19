package build

import "github.com/arschles/gci/tpl"

type Pipeline struct {
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"` // can be templated
	Steps       []PipelineStep `yaml:"steps"`
}

func (p Pipeline) RenderDescription(varMap VarMap) (string, error) {
	return tpl.Render(p.Name, p.Description, varMap)
}

type PipelineStep struct {
	Name   string              `yaml:"name"` // can be templated
	Params []PipelineStepParam `yaml:"params"`
}
type PipelineStepParam struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"` // can be templated
}

func (p PipelineStepParam) RenderValue(varMap VarMap) (string, error) {
	return tpl.Render(p.Name, p.Value, varMap)
}
