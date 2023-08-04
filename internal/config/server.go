package config

import (
	"flag"
)

type ServerConfig struct {
	Addr string
}

func InitServerConfig() *ServerConfig {
	conf := new(ServerConfig)
	addr := flag.String("a", "localhost:8080", "Server address. Usage: -a=host:port")
	flag.Parse()
	conf.Addr = *addr
	return conf
}
