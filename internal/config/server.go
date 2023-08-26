package config

import (
	"flag"

	"github.com/caarlos0/env/v9"
)

const (
	ServerDefAddr = "localhost:8080"
	templatePath  = "web/template/"
)

type ServerConfig struct {
	Addr         string `env:"ADDRESS"`
	TemplatePath string
}

func InitServerConfig() (*ServerConfig, error) {
	conf := new(ServerConfig)
	err := env.Parse(conf)
	if err != nil {
		return nil, err
	}
	conf.TemplatePath = templatePath
	if conf.Addr != "" {
		return conf, nil
	}
	flag.StringVar(&conf.Addr, "a", "localhost:8080", "Server address. Usage: -a=host:port")
	flag.Parse()
	return conf, nil
}
