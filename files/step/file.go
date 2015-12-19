package step

type File struct {
	Version string   `yaml:"version"`
	Image   string   `yaml:"image"`
	Command string   `yaml:"command"`
	Params  []Param  `yaml:"params"`
	Envs    []Env    `yaml:"envs"`
	Volumes []Volume `yaml:"volumes"`
}
