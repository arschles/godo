package ci

const (
	DefaultBindHost   = "0.0.0.0"
	DefaultClientHost = "127.0.0.1"
	DefaultPort       = 8083
)

type CI struct {
	Server Server `yaml:"server"`
	Build  Build  `yaml:"build"`
}

type Server struct {
	Host string `yaml:"host"`
	Port uint   `yaml:"port"`
}

func (c Server) GetBindHost() string {
	if c.Host == "" {
		return DefaultBindHost
	}
	return c.Host
}

func (c Server) GetPort() uint {
	if c.Port == 0 {
		return DefaultPort
	}
	return c.Port
}
