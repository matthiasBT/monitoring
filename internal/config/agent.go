package config

import (
	"flag"
	"github.com/caarlos0/env/v9"
	"log"
)

type AgentConfig struct {
	ServerAddr     string `env:"ADDRESS"`
	ReportInterval uint   `env:"REPORT_INTERVAL"`
	PollInterval   uint   `env:"POLL_INTERVAL"`
}

func InitAgentConfig() *AgentConfig {
	conf := new(AgentConfig)
	err := env.Parse(conf)
	if err != nil {
		log.Fatal(err)
	}
	addr := flag.String("a", "localhost:8080", "Server address. Usage: -a=host:port")
	reportInterval := flag.Uint("r", 10, "How often to send metrics to the server, seconds")
	pollInterval := flag.Uint("p", 2, "How often to query metrics, seconds")
	flag.Parse()
	if conf.ServerAddr == "" {
		conf.ServerAddr = *addr
	}
	if conf.ReportInterval == 0 {
		conf.ReportInterval = *reportInterval
	}
	if conf.PollInterval == 0 {
		conf.PollInterval = *pollInterval
	}
	return conf
}
