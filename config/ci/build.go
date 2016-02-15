package ci

type Build struct {
	Excludes     []Exclude `yaml:"exclude"`
	Host         string    `yaml:"host"`
	Port         uint      `yaml:"port"`
	CrossCompile bool      `yaml:"cross-compile"`
	Env          []string  `yaml:"env"`
}

func (b Build) GetHost() string {
	if b.Host == "" {
		return DefaultClientHost
	}
	return b.Host
}

func (b Build) GetPort() uint {
	if b.Port == 0 {
		return DefaultPort
	}
	return b.Port
}

type Exclude struct {
	Name      string `yaml:"name"`
	Recursive bool   `yaml:"recursive"`
}
