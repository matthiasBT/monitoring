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
	DefRateLimit            = 1
)

type Config struct {
	Addr                 string `env:"ADDRESS"`
	UpdateURL            string
	ReportInterval       uint   `env:"REPORT_INTERVAL"`
	PollInterval         uint   `env:"POLL_INTERVAL"`
	HMACKey              string `env:"KEY"`
	RateLimit            uint   `env:"RATE_LIMIT"`
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
	hmacKey := flag.String("k", "", "HMAC key for integrity checks")
	rateLimit := flag.Uint("l", DefRateLimit, "Max number of active workers")
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
	if conf.HMACKey == "" {
		conf.HMACKey = *hmacKey
	}
	if conf.RateLimit == 0 {
		conf.RateLimit = *rateLimit
	}
	conf.UpdateURL = updateURL
	conf.RetryAttempts = DefRetryAttempts
	conf.RetryIntervalInitial = DefRetryIntervalInitial
	conf.RetryIntervalBackoff = DefRetryIntervalBackoff
	return conf, nil
}
