// Package agent provides the configuration setup for an agent application.
// It includes structures and functions for initializing and managing
// configuration settings.
package agent

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
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
	// ConfigPath is the path of a JSON configuration file
	ConfigPath string `env:"CONFIG"`

	// Addr represents the server address to which the agent connects.
	Addr string `env:"ADDRESS" json:"address"`

	// UpdateURL is the URL endpoint for sending updates.
	UpdateURL string

	// HMACKey is used for HMAC-based integrity checks.
	HMACKey string `env:"KEY"`

	// CryptoKey is the public key of the monitoring server
	CryptoKey string `env:"CRYPTO_KEY" json:"crypto_key"`

	// ReportInterval specifies how often (in seconds) the agent sends metrics to the server.
	ReportInterval uint `env:"REPORT_INTERVAL" json:"report_interval"`

	// PollInterval specifies how often (in seconds) the agent queries for metrics.
	PollInterval uint `env:"POLL_INTERVAL" json:"poll_interval"`

	// RateLimit defines the maximum number of active workers for processing.
	RateLimit uint `env:"RATE_LIMIT"`

	// RetryAttempts is the number of retry attempts for failed requests.
	RetryAttempts int

	// RetryIntervalInitial is the initial time duration between retries.
	RetryIntervalInitial time.Duration

	// RetryIntervalBackoff is the time duration for exponential backoff between retries.
	RetryIntervalBackoff time.Duration
}

// ReadServerPublicKey reads a file and return an RSA private key
func (c *Config) ReadServerPublicKey() (*rsa.PublicKey, error) {
	if c.CryptoKey == "" {
		return nil, nil
	}

	data, err := os.ReadFile(c.CryptoKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, fmt.Errorf("unknown type of public key")
	}
}

// InitConfig initializes the Config structure by parsing environment variables
// and command-line flags. It provides defaults for missing values and sets up
// the configuration for the agent.
func InitConfig() (*Config, error) {
	conf := new(Config)
	flag.StringVar(&conf.ConfigPath, "c", "", "Configuration file path")
	flag.StringVar(&conf.Addr, "a", DefAddr, "Server address. Usage: -a=host:port")
	flag.UintVar(
		&conf.ReportInterval, "r", DefReportInterval, "How often to send metrics to the server, seconds",
	)
	flag.UintVar(&conf.PollInterval, "p", DefPollInterval, "How often to query metrics, seconds")
	flag.StringVar(&conf.HMACKey, "k", "", "HMAC key for integrity checks")
	flag.StringVar(&conf.CryptoKey, "crypto-key", "", "Path to a file with the server public key")
	flag.UintVar(&conf.RateLimit, "l", DefRateLimit, "Max number of active workers")
	flag.Parse()
	if jsonConfigPath, ok := os.LookupEnv("CONFIG"); ok {
		conf.ConfigPath = jsonConfigPath
	}
	if conf.ConfigPath != "" {
		raw, err := os.ReadFile(conf.ConfigPath)
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(raw, conf); err != nil {
			panic(err)
		}
		flag.Parse()
	}
	err := env.Parse(conf)
	if err != nil {
		return nil, err
	}
	conf.UpdateURL = updateURL
	conf.RetryAttempts = DefRetryAttempts
	conf.RetryIntervalInitial = DefRetryIntervalInitial
	conf.RetryIntervalBackoff = DefRetryIntervalBackoff
	return conf, nil
}
