// Package server contains the configuration and initialization logic for a server application.
// It defines structures and functions for setting up server configuration including
// storage, restoration, and retry logic.
package server

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v9"
)

const (
	templatePath            = "web/template/"
	DefAddr                 = "localhost:8080"
	DefStoreInterval        = 300
	DefFileStoragePath      = "/tmp/metrics-db.json"
	DefRestore              = true
	DefRetryAttempts        = 3
	DefRetryIntervalInitial = 1 * time.Second
	DefRetryIntervalBackoff = 2 * time.Second
)

// Config defines the configuration parameters for the server. It includes server address,
// template path, storage settings, HMAC key, and retry settings.
type Config struct {
	// Addr represents the server address and port.
	Addr string `env:"ADDRESS"`

	// TemplatePath is the file path to the web templates used by the server.
	TemplatePath string

	// StoreInterval specifies the interval (in seconds) for storing data to the file.
	StoreInterval *uint `env:"STORE_INTERVAL"`

	// FileStoragePath is the file path for storing metrics data.
	FileStoragePath string `env:"FILE_STORAGE_PATH"`

	// Restore indicates whether to restore the initial state from the file.
	Restore *bool `env:"RESTORE"`

	// DatabaseDSN is the Data Source Name for connecting to a PostgreSQL database.
	DatabaseDSN string `env:"DATABASE_DSN"`

	// HMACKey is used for HMAC-based integrity checks.
	HMACKey string `env:"KEY"`

	// RetryAttempts is the number of retry attempts for failed requests.
	RetryAttempts int

	// RetryIntervalInitial is the initial duration between retries.
	RetryIntervalInitial time.Duration

	// RetryIntervalBackoff is the duration for exponential backoff between retries.
	RetryIntervalBackoff time.Duration
}

// InitConfig initializes the Config structure by parsing environment variables
// and command-line flags. It sets defaults for missing values and prepares
// the server configuration.
func InitConfig() (*Config, error) {
	conf := new(Config)
	err := env.Parse(conf)
	if err != nil {
		return nil, err
	}

	conf.TemplatePath = templatePath
	conf.RetryAttempts = DefRetryAttempts
	conf.RetryIntervalInitial = DefRetryIntervalInitial
	conf.RetryIntervalBackoff = DefRetryIntervalBackoff

	flagAddr := flag.String("a", DefAddr, "Server address. Usage: -a=host:port")

	var flagStoragePath string
	flag.StringVar(&flagStoragePath, "f", DefFileStoragePath, "Path to storage file")

	var flagDatabaseDSN string
	flag.StringVar(&flagDatabaseDSN, "d", "", "PostgreSQL database DSN")

	flagRestore := flag.Bool("r", DefRestore, "Restore init state from the file (see -f flag)")
	flagStoreInterval := flag.Uint("i", DefStoreInterval, "How often to store data in the file")
	hmacKey := flag.String("k", "", "HMAC key for integrity checks")
	flag.Parse()

	if conf.Addr == "" {
		conf.Addr = *flagAddr
	}
	if conf.FileStoragePath == "" {
		conf.FileStoragePath = flagStoragePath
	}
	if conf.DatabaseDSN == "" {
		conf.DatabaseDSN = flagDatabaseDSN
	}
	if conf.Restore == nil {
		conf.Restore = flagRestore
	}
	if conf.StoreInterval == nil {
		conf.StoreInterval = flagStoreInterval
	}
	if conf.HMACKey == "" {
		conf.HMACKey = *hmacKey
	}
	return conf, nil
}

// FlushesSync determines if the server is configured to flush data synchronously.
// Returns true if the StoreInterval is set to 0, indicating synchronous flush.
func (c *Config) FlushesSync() bool {
	return *c.StoreInterval == 0
}

// Flushes checks whether the server is configured to flush data to storage.
// Returns true if either FileStoragePath or DatabaseDSN is configured.
func (c *Config) Flushes() bool {
	return c.FileStoragePath != "" || c.DatabaseDSN != ""
}
