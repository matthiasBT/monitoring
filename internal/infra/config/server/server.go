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

type Config struct {
	Addr                 string `env:"ADDRESS"`
	TemplatePath         string
	StoreInterval        *uint  `env:"STORE_INTERVAL"`
	FileStoragePath      string `env:"FILE_STORAGE_PATH"`
	Restore              *bool  `env:"RESTORE"`
	DatabaseDSN          string `env:"DATABASE_DSN"`
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
	return conf, nil
}

func (c *Config) FlushesSync() bool {
	return *c.StoreInterval == 0
}

func (c *Config) Flushes() bool {
	return c.FileStoragePath != ""
}
