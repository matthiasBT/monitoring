// Package main is the entry point for the server application in the monitoring system.
// It sets up and runs the HTTP server, handling configuration, logging, and graceful shutdowns.
// The package integrates various internal components like routing, data storage, and middleware.
package main

import (
	"context"
	"crypto/rsa"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/matthiasBT/monitoring/internal/infra/compression"
	"github.com/matthiasBT/monitoring/internal/infra/config/server"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/secure"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
	"github.com/matthiasBT/monitoring/internal/server/adapters"
	"github.com/matthiasBT/monitoring/internal/server/entities"
	"github.com/matthiasBT/monitoring/internal/server/usecases"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// setupServer configures and returns a new HTTP router with middleware and routes.
// It includes logging, compression, optional HMAC checking, and controller routes.
func setupServer(
	logger logging.ILogger, controller *usecases.BaseController, hmacKey string, key *rsa.PrivateKey,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(logging.Middleware(logger))
	r.Use(compression.MiddlewareReader, compression.MiddlewareWriter)
	if hmacKey != "" {
		r.Use(secure.MiddlewareHashReader(hmacKey), secure.MiddlewareHashWriter(hmacKey))
	}
	if key != nil {
		r.Use(secure.MiddlewareCryptoReader(key))
	}
	r.Mount("/", controller.Route())
	return r
}

// gracefulShutdown handles the graceful shutdown of the server.
// It listens for system signals and shuts down the server after processing ongoing requests.
func gracefulShutdown(srv *http.Server, done chan struct{}, logger logging.ILogger) {
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quitChannel
	logger.Infof("Received signal: %v\n", sig)
	done <- struct{}{}
	time.Sleep(2 * time.Second)

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err.Error())
	}
}

// setupRetrier configures and returns a Retrier based on the provided server configuration.
// It sets up retry attempts, intervals, and logging for handling network-related retries.
func setupRetrier(conf *server.Config, logger logging.ILogger) utils.Retrier {
	return utils.Retrier{
		Attempts:         conf.RetryAttempts,
		IntervalFirst:    conf.RetryIntervalInitial,
		IntervalIncrease: conf.RetryIntervalBackoff,
		Logger:           logger,
	}
}

// setupKeeper initializes and returns the appropriate Keeper (database or file) based on configuration.
// It configures the storage mechanism for the server, handling data persistence.
func setupKeeper(conf *server.Config, logger logging.ILogger, retrier utils.Retrier) entities.Keeper {
	if conf.Flushes() {
		if conf.DatabaseDSN != "" {
			return adapters.NewDBKeeper(conf, logger, retrier)
		} else {
			return adapters.NewFileKeeper(conf, logger, retrier)
		}
	}
	return nil
}

// setupTicker creates and returns a ticker channel based on the configuration.
// It's used for periodic operations like data flushing.
func setupTicker(conf *server.Config) <-chan time.Time {
	if conf.FlushesSync() {
		return make(chan time.Time) // will never be used
	} else {
		ticker := time.NewTicker(time.Duration(conf.StoreInterval) * time.Second)
		return ticker.C
	}
}

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
	tickerChan := setupTicker(conf)
	retrier := setupRetrier(conf, logger)

	keeper := setupKeeper(conf, logger, retrier)
	if keeper != nil {
		defer keeper.Shutdown()
	}
	storage := adapters.NewMemStorage(done, tickerChan, logger, keeper)

	if conf.Flushes() {
		if conf.Restore {
			state := keeper.Restore()
			storage.Init(state)
		}
		if !conf.FlushesSync() {
			go storage.FlushPeriodic(context.Background())
		}
	}

	controller := usecases.NewBaseController(logger, storage, conf.TemplatePath)
	key, err := conf.ReadPrivateKey()
	if err != nil {
		panic(err)
	}
	r := setupServer(logger, controller, conf.HMACKey, key)
	srv := http.Server{Addr: conf.Addr, Handler: r}
	go func() {
		logger.Infof("Launching the server at %s\n", conf.Addr)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err)
		}
	}()

	gracefulShutdown(&srv, done, logger)
}
