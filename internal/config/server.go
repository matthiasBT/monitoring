package config

import (
	"flag"
	"github.com/caarlos0/env/v9"
	"log"
)

type ServerConfig struct {
	Addr string `env:"ADDRESS"`
}

func InitServerConfig() *ServerConfig {
	conf := new(ServerConfig)
	err := env.Parse(conf)
	if err != nil {
		log.Fatal(err)
	}
	if conf.Addr != "" {
		return conf
	}
	flag.StringVar(&conf.Addr, "a", "localhost:8080", "Server address. Usage: -a=host:port")
	flag.Parse()
	return conf
}
