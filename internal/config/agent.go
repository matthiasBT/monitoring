package config

import (
	"flag"
	"time"
)

type AgentConfig struct {
	ServerAddr     string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func InitAgentConfig() *AgentConfig {
	conf := new(AgentConfig)
	addr := flag.String("a", "localhost:8080", "Server address. Usage: -a=host:port")
	reportInterval := flag.Duration("r", 10, "How often to send metrics to the server, seconds")
	pollInterval := flag.Duration("p", 2, "How often to query metrics, seconds")
	flag.Parse()
	conf.ServerAddr = *addr
	conf.ReportInterval = *reportInterval * time.Second
	conf.PollInterval = *pollInterval * time.Second
	return conf
}
