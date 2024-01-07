// Package server contains the configuration and initialization logic for a server application.
// It defines structures and functions for setting up server configuration including
// storage, restoration, and retry logic.
package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"os"
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
	// ConfigPath is the path of a JSON configuration file
	ConfigPath string `env:"CONFIG"`

	// Addr represents the server address and port.
	Addr string `env:"ADDRESS" json:"address"`

	// TemplatePath is the file path to the web templates used by the server.
	TemplatePath string

	// StoreInterval specifies the interval (in seconds) for storing data to the file.
	StoreInterval uint `env:"STORE_INTERVAL" json:"store_interval"`

	// FileStoragePath is the file path for storing metrics data.
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"store_file"`

	// Restore indicates whether to restore the initial state from the file.
	Restore bool `env:"RESTORE" json:"restore"`

	// DatabaseDSN is the Data Source Name for connecting to a PostgreSQL database.
	DatabaseDSN string `env:"DATABASE_DSN" json:"database_dsn"`

	// HMACKey is used for HMAC-based integrity checks.
	HMACKey string `env:"KEY" json:"key"`

	// CryptoKey is used for payload decryption
	CryptoKey string `env:"CRYPTO_KEY" json:"crypto_key"`

	// RetryAttempts is the number of retry attempts for failed requests.
	RetryAttempts int

	// RetryIntervalInitial is the initial duration between retries.
	RetryIntervalInitial time.Duration

	// RetryIntervalBackoff is the duration for exponential backoff between retries.
	RetryIntervalBackoff time.Duration
}

// ReadPrivateKey reads a file and returns an RSA private key
func (c *Config) ReadPrivateKey() (*rsa.PrivateKey, error) {
	if c.CryptoKey == "" {
		return nil, nil
	}

	data, err := os.ReadFile(c.CryptoKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key.(*rsa.PrivateKey), nil
}

// InitConfig initializes the Config structure by parsing environment variables
// and command-line flags. It sets defaults for missing values and prepares
// the server configuration.
func InitConfig() (*Config, error) {
	conf := new(Config)
	flag.StringVar(&conf.ConfigPath, "c", "", "Configuration file path")
	flag.StringVar(&conf.Addr, "a", DefAddr, "Server address. Usage: -a=host:port")
	flag.StringVar(&conf.FileStoragePath, "f", DefFileStoragePath, "Path to storage file")
	flag.StringVar(&conf.DatabaseDSN, "d", "", "PostgreSQL database DSN")
	flag.BoolVar(&conf.Restore, "r", DefRestore, "Restore init state from the file (see -f flag)")
	flag.UintVar(&conf.StoreInterval, "i", DefStoreInterval, "How often to store data in the file")
	flag.StringVar(&conf.HMACKey, "k", "", "HMAC key for integrity checks")
	flag.StringVar(&conf.CryptoKey, "crypto-key", "", "Path to a file with the server private key")
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
	conf.TemplatePath = templatePath
	conf.RetryAttempts = DefRetryAttempts
	conf.RetryIntervalInitial = DefRetryIntervalInitial
	conf.RetryIntervalBackoff = DefRetryIntervalBackoff
	return conf, nil
}

// FlushesSync determines if the server is configured to flush data synchronously.
// Returns true if the StoreInterval is set to 0, indicating synchronous flush.
func (c *Config) FlushesSync() bool {
	return c.StoreInterval == 0
}

// Flushes checks whether the server is configured to flush data to storage.
// Returns true if either FileStoragePath or DatabaseDSN is configured.
func (c *Config) Flushes() bool {
	return c.FileStoragePath != "" || c.DatabaseDSN != ""
}
