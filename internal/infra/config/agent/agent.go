package agent

import (
	"flag"

	"github.com/caarlos0/env/v9"
)

const (
	updateURL         = "/updates/"
	DefAddr           = "localhost:8080"
	DefReportInterval = 10
	DefPollInterval   = 2
)

type Config struct {
	Addr           string `env:"ADDRESS"`
	UpdateURL      string
	ReportInterval uint `env:"REPORT_INTERVAL"`
	PollInterval   uint `env:"POLL_INTERVAL"`
}

func InitConfig() (*Config, error) {
	conf := new(Config)
	err := env.Parse(conf)
	if err != nil {
		return nil, err
	}
	addr := flag.String("a", DefAddr, "Server address. Usage: -a=host:port")
	reportInterval := flag.Uint(
		"r", DefReportInterval, "How often to send metrics to the server, seconds",
	)
	pollInterval := flag.Uint("p", DefPollInterval, "How often to query metrics, seconds")
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
	conf.UpdateURL = updateURL
	return conf, nil
}
