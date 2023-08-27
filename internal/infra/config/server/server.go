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
	var restore bool
	var storeInterval uint
	conf.TemplatePath = templatePath
	if conf.Addr == "" {
		flag.StringVar(&conf.Addr, "a", "localhost:8080", "Server address. Usage: -a=host:port")
	}
	if conf.FileStoragePath == "" {
		flag.StringVar(&conf.FileStoragePath, "f", DefFileStoragePath, "Path to storage file")
	}
	if conf.Restore == nil {
		flag.BoolVar(&restore, "r", DefRestore, "Restore init state from the file (see -f flag)")
	}
	if conf.StoreInterval == nil {
		flag.UintVar(&storeInterval, "i", DefStoreInterval, "How often to store data in the file")
	}
	flag.Parse()
	conf.Restore = &restore
	conf.StoreInterval = &storeInterval
	return conf, nil
}
