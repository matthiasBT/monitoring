package agent

import (
	"flag"

	"github.com/caarlos0/env/v9"
)

const (
	AgentDefAddr           = "localhost:8080"
	AgentDefReportInterval = 10
	AgentDefPollInterval   = 2
)

type Config struct {
	Addr           string `env:"ADDRESS"`
	ReportInterval uint   `env:"REPORT_INTERVAL"`
	PollInterval   uint   `env:"POLL_INTERVAL"`
}

func InitConfig() (*Config, error) {
	conf := new(Config)
	err := env.Parse(conf)
	if err != nil {
		return nil, err
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
	return conf, nil
}
