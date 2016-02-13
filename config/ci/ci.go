package ci

const (
	DefaultHost = "0.0.0.0"
	DefaultPort = 8083
)

type CI struct {
	Server Server `yaml:"server"`
}

type Server struct {
	Host string `yaml:"host"`
	Port uint   `yaml:"port"`
}

func (c Server) GetHost() string {
	if c.Host == "" {
		return DefaultHost
	}
	return c.Host
}
func (c Server) GetPort() uint {
	if c.Port == 0 {
		return DefaultPort
	}
	return c.Port
}
