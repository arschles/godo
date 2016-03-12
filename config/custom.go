package config

type CustomTarget struct {
	ImageName string        `yaml:"image-name"`
	ImageTag  string        `yaml:"image-tag"`
	Command   string        `yaml:"command"`
	Mounts    []CustomMount `yaml:"mounts"`
}

type CustomMount struct {
	Host      string `yaml:"host"`
	Container string `yaml:"container"`
}
