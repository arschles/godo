package main

type Consfile struct {
	Bootstrap     *Bootstrap     `yaml:"bootstrap"`
	Build         *Build         `yaml:"build"`
	Test          *Test          `yaml:"test"`
	Install       *Install       `yaml:"install"`
	OtherCommands []OtherCommand `yaml:"others"`
}

type Bootstrap struct {
	Commands []string `yaml:"commands"`
}

type Build struct {
	Depends string   `yaml:"depends"`
	Output  string   `yaml:"output"`
	Env     []string `yaml:"env"`
}

type Install struct {
	Depends    string   `yaml:"depends"`
	PreScript  string   `yaml:"pre_script"`
	PostScript string   `yaml:"post_script"`
	Env        []string `yaml:"env"`
}

type Test struct {
	Depends string   `yaml:"depends"`
	Verbose bool     `yaml:"verbose"`
	Paths   []string `yaml:"paths"`
}

type OtherCommand struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Command     string `yaml:"cmd"`
	Depends     string `yaml:"depends"`
}
