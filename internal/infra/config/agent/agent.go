package agent

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v9"
)

const (
	updateURL               = "/updates/"
	DefAddr                 = "localhost:8080"
	DefReportInterval       = 10
	DefPollInterval         = 2
	DefRetryAttempts        = 3
	DefRetryIntervalInitial = 1 * time.Second
	DefRetryIntervalBackoff = 2 * time.Second
)

type Config struct {
	Addr                 string `env:"ADDRESS"`
	UpdateURL            string
	ReportInterval       uint `env:"REPORT_INTERVAL"`
	PollInterval         uint `env:"POLL_INTERVAL"`
	RetryAttempts        int
	RetryIntervalInitial time.Duration
	RetryIntervalBackoff time.Duration
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
	conf.RetryAttempts = DefRetryAttempts
	conf.RetryIntervalInitial = DefRetryIntervalInitial
	conf.RetryIntervalBackoff = DefRetryIntervalBackoff
	return conf, nil
}
