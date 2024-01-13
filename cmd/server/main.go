// Package main is the entry point for the server application in the monitoring system.
// It sets up and runs the HTTP server, handling configuration, logging, and graceful shutdowns.
// The package integrates various internal components like routing, data storage, and middleware.
package main

import (
	"errors"
	"net/http"

	"github.com/matthiasBT/monitoring/internal/infra/config/server"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
	"github.com/matthiasBT/monitoring/internal/server/adapters"
	"github.com/matthiasBT/monitoring/internal/server/startup"
	"github.com/matthiasBT/monitoring/internal/server/usecases"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// main is the entry function of the application. It sets up and starts the HTTP server,
// including configuration, logging, storage, and routing. It also manages the application's lifecycle,
// handling initialization and graceful shutdown.
func main() {
	utils.PrintBuildFlags(buildVersion, buildDate, buildCommit)

	logger := logging.SetupLogger()
	conf, err := server.InitConfig()
	if err != nil {
		logger.Fatal(err)
	}

	done := make(chan struct{}, 1)
	tickerChan := startup.SetupTicker(conf)
	retrier := startup.SetupRetrier(conf, logger)

	keeper := startup.SetupKeeper(conf, logger, retrier)
	if keeper != nil {
		defer keeper.Shutdown()
	}
	storage := adapters.NewMemStorage(done, tickerChan, logger, keeper)
	startup.PrepareStorage(conf, keeper, storage)
	controller := usecases.NewBaseController(logger, storage, conf.TemplatePath)
	key, err := conf.ReadPrivateKey()
	if err != nil {
		panic(err)
	}
	r := startup.SetupHTTPServer(logger, controller, conf.HMACKey, conf.TrustedSubnet, key)
	srv := http.Server{Addr: conf.Addr, Handler: r}
	go func() {
		logger.Infof("Launching the server at %s\n", conf.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err)
		}
	}()

	startup.GracefulShutdownHTTP(&srv, done, logger)
}
