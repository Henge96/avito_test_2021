package config

type Config struct {
	Server
	Db
}

type Server struct {
	Host string `yaml:"host"`
	Port
}

type Db struct {
	Driver   string `yaml:"driver"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	HostDb   string `yaml:"hostdb"`
	PortDb   string `yaml:"portdb"`
	Dbname   string `yaml:"dbname"`
	Mode     string `yaml:"mode"`
}

type Port struct {
	Http string `yaml:"http"`
}
