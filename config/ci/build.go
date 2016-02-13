package ci

type Build struct {
	Exclude      []string `yaml:"exclude"`
	Host         string   `yaml:"host"`
	Port         uint     `yaml:"port"`
	CrossCompile bool     `yaml:"cross-compile"`
	Env          []string `yaml:"env"`
}

func (b Build) GetHost() string {
	if b.Host == "" {
		return DefaultHost
	}
	return b.Host
}

func (b Build) GetPort() uint {
	if b.Port == 0 {
		return DefaultPort
	}
	return b.Port
}
