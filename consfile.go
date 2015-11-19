package main

type Consfile struct {
	Bootstrap Bootstrap `yaml:"bootstrap"`
	Build     Build     `yaml:"build"`
	Test      Test      `yaml:"test"`
	Install   Install   `yaml:"install"`
}

type Bootstrap struct {
	Command string `yaml:"output"`
}

type Build struct {
	Depends string `yaml:"depends"`
	Output  string `yaml:"output"`
}

type Test struct {
	Depends string   `yaml:"depends"`
	Verbose bool     `yaml:"verbose"`
	Paths   []string `yaml:"paths"`
}

type Other struct {
	Name    string `yaml:"name"`
	Command string `yaml:"cmd"`
	Depends string `yaml:"depends"`
}

type Install struct {
	//TODO
}
