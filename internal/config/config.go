package config

import (
	"flag"

	"github.com/caarlos0/env/v9"
	"github.com/matthiasBT/monitoring/internal/interfaces"
)

const (
	ServerDefAddr          = "localhost:8080"
	AgentDefAddr           = "localhost:8080"
	AgentDefReportInterval = 10
	AgentDefPollInterval   = 2
	templatePath           = "web/template/"
)

type AgentConfig struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval uint   `env:"REPORT_INTERVAL"`
	PollInterval   uint   `env:"POLL_INTERVAL"`
}

type ServerConfig struct {
	Addr         string `env:"ADDRESS"`
	TemplatePath string
}

func InitAgentConfig(logger interfaces.ILogger) *AgentConfig {
	conf := new(AgentConfig)
	err := env.Parse(conf)
	if err != nil {
		logger.Fatal(err)
	}
	addr := flag.String("a", AgentDefAddr, "Server address. Usage: -a=host:port")
	reportInterval := flag.Uint(
		"r", AgentDefReportInterval, "How often to send metrics to the server, seconds",
	)
	pollInterval := flag.Uint("p", AgentDefPollInterval, "How often to query metrics, seconds")
	flag.Parse()
	if conf.Addr == "" {
		conf.Addr = *addr
	}
	if conf.ReportInterval == 0 {
		conf.ReportInterval = *reportInterval
	}
	if conf.PollInterval == 0 {
		conf.PollInterval = *pollInterval
	}
	fmt.Printf("Agent config: %v\n", conf) // todo: remove
	return conf
}

func InitServerConfig(logger interfaces.ILogger) *ServerConfig {
	conf := new(ServerConfig)
	err := env.Parse(conf)
	if err != nil {
		logger.Fatal(err)
	}
	conf.TemplatePath = templatePath
	if conf.Addr != "" {
		return conf
	}
	flag.StringVar(&conf.Addr, "a", "localhost:8080", "Server address. Usage: -a=host:port")
	flag.Parse()
	fmt.Printf("Server config: %v\n", conf) // todo: remove
	return conf
}
