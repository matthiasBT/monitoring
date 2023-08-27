package server

import (
	"flag"

	"github.com/caarlos0/env/v9"
)

const (
	templatePath       = "web/template/"
	DefAddr            = "localhost:8080"
	DefStoreInterval   = 300
	DefFileStoragePath = "/tmp/metrics-db.json"
	DefRestore         = true
)

type Config struct {
	Addr            string `env:"ADDRESS"`
	TemplatePath    string
	StoreInterval   *uint  `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         *bool  `env:"RESTORE"`
}

func InitConfig() (*Config, error) {
	conf := new(Config)
	err := env.Parse(conf)
	if err != nil {
		return nil, err
	}
	conf.TemplatePath = templatePath
	flagAddr := flag.String("a", DefAddr, "Server address. Usage: -a=host:port")
	flagStoragePath := flag.String("f", DefFileStoragePath, "Path to storage file")
	flagRestore := flag.Bool("r", DefRestore, "Restore init state from the file (see -f flag)")
	flagStoreInterval := flag.Uint("i", DefStoreInterval, "How often to store data in the file")
	flag.Parse()
	if conf.Addr == "" {
		conf.Addr = *flagAddr
	}
	if conf.FileStoragePath == "" {
		conf.FileStoragePath = *flagStoragePath
	}
	if conf.Restore == nil {
		conf.Restore = flagRestore
	}
	if conf.StoreInterval == nil {
		conf.StoreInterval = flagStoreInterval
	}
	return conf, nil
}

func (c *Config) StoresSync() bool {
	return *c.StoreInterval == 0
}
