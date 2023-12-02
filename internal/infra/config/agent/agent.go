// Package agent provides the configuration setup for an agent application.
// It includes structures and functions for initializing and managing
// configuration settings.

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

// Config defines the configuration parameters for the agent. It includes
// server address, update URL, intervals for reporting and polling metrics,
// HMAC key for integrity checks, rate limits, and retry settings.
type Config struct {
	// Addr represents the server address to which the agent connects.
	Addr string `env:"ADDRESS"`

	// UpdateURL is the URL endpoint for sending updates.
	UpdateURL string

	// HMACKey is used for HMAC-based integrity checks.
	HMACKey string `env:"KEY"`

	// ReportInterval specifies how often (in seconds) the agent sends metrics to the server.
	ReportInterval uint `env:"REPORT_INTERVAL"`

	// PollInterval specifies how often (in seconds) the agent queries for metrics.
	PollInterval uint `env:"POLL_INTERVAL"`

	// RateLimit defines the maximum number of active workers for processing.
	RateLimit uint `env:"RATE_LIMIT"`

	// RetryAttempts is the number of retry attempts for failed requests.
	RetryAttempts int

	// RetryIntervalInitial is the initial time duration between retries.
	RetryIntervalInitial time.Duration

	// RetryIntervalBackoff is the time duration for exponential backoff between retries.
	RetryIntervalBackoff time.Duration
}

// InitConfig initializes the Config structure by parsing environment variables
// and command-line flags. It provides defaults for missing values and sets up
// the configuration for the agent.
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
