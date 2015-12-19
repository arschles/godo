package build

import (
	"errors"

	"github.com/arschles/gci/tpl"
)

var (
	ErrNoName        = errors.New("no name")
	ErrNoValue       = errors.New("no value")
	ErrNoDescription = errors.New("no description")
	ErrNoSteps       = errors.New("no steps")
)

type Pipeline struct {
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"` // can be templated
	Steps       []PipelineStep `yaml:"steps"`
}

func (p Pipeline) Validate() error {
	if p.Name == "" {
		return ErrNoName
	}
	if p.Description == "" {
		return ErrNoDescription
	}
	if len(p.Steps) == 0 {
		return ErrNoSteps
	}
	return nil
}

func (p Pipeline) RenderDescription(varMap VarMap) (string, error) {
	return tpl.Render(p.Name, p.Description, varMap)
}

type PipelineStep struct {
	Name   string              `yaml:"name"` // can be templated
	Params []PipelineStepParam `yaml:"params"`
}

func (p PipelineStep) Validate() error {
	if p.Name == "" {
		return ErrNoName
	}
	return nil
}

type PipelineStepParam struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"` // can be templated
}

func (p PipelineStepParam) Validate() error {
	if p.Name == "" {
		return ErrNoName
	}
	if p.Value == "" {
		return ErrNoValue
	}
	return nil
}

func (p PipelineStepParam) RenderValue(varMap VarMap) (string, error) {
	return tpl.Render(p.Name, p.Value, varMap)
}
