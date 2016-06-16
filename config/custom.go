package config

// CustomTarget is a custom command that executes a set of commands in a container
type CustomTarget struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	ImageName   string   `yaml:"image-name"`
	ImageTag    string   `yaml:"image-tag"`
	Command     string   `yaml:"command"`
	MountTarget string   `yaml:"mount-target"`
	Envs        []string `yaml:"environment"`
}
