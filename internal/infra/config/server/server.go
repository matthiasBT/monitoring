package server

import (
	"flag"

	"github.com/caarlos0/env/v9"
)

const (
	ServerDefAddr = "localhost:8080"
	templatePath  = "web/template/"
)

type Config struct {
	Addr         string `env:"ADDRESS"`
	TemplatePath string
}

func InitConfig() (*Config, error) {
	conf := new(Config)
	err := env.Parse(conf)
	if err != nil {
		return nil, err
	}
	conf.TemplatePath = templatePath
	if conf.Addr == "" {
		flag.StringVar(&conf.Addr, "a", "localhost:8080", "Server address. Usage: -a=host:port")
		flag.Parse()
	}
	return conf, nil
}
