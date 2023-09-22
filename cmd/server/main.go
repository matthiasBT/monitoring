package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/matthiasBT/monitoring/internal/infra/utils"
	"github.com/matthiasBT/monitoring/internal/server/entities"

	"github.com/go-chi/chi/v5"
	"github.com/matthiasBT/monitoring/internal/infra/compression"
	"github.com/matthiasBT/monitoring/internal/infra/config/server"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/adapters"
	"github.com/matthiasBT/monitoring/internal/server/usecases"
)

func setupServer(logger logging.ILogger, controller *usecases.BaseController) *chi.Mux {
	r := chi.NewRouter()
	r.Use(logging.Middleware(logger))
	r.Use(compression.MiddlewareReader)
	r.Use(compression.MiddlewareWriter)
	r.Mount("/", controller.Route())
	return r
}

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

func setupRetrier(conf *server.Config, logger logging.ILogger) utils.Retrier {
	return utils.Retrier{
		Attempts:         conf.RetryAttempts,
		IntervalFirst:    conf.RetryIntervalInitial,
		IntervalIncrease: conf.RetryIntervalBackoff,
		Logger:           logger,
	}
}

func setupTicker(conf *server.Config) <-chan time.Time {
	if conf.FlushesSync() {
		return make(chan time.Time) // will never be used
	} else {
		ticker := time.NewTicker(time.Duration(*conf.StoreInterval) * time.Second)
		return ticker.C
	}
}

func main() {
	logger := logging.SetupLogger()
	conf, err := server.InitConfig()
	if err != nil {
		logger.Fatal(err)
	}

	done := make(chan struct{}, 1)
	tickerChan := setupTicker(conf)
	retrier := setupRetrier(conf, logger)

	storage := adapters.NewMemStorage(done, tickerChan, logger, nil)
	var keeper entities.Keeper

	if conf.Flushes() {
		if conf.DatabaseDSN != "" {
			keeper = adapters.NewDBKeeper(conf, logger, retrier)
		} else {
			keeper = adapters.NewFileKeeper(conf, logger, retrier)
		}
		defer keeper.Shutdown()
		if *conf.Restore {
			state := keeper.Restore()
			storage.Init(state)
		}
		storage.SetKeeper(keeper)
		if !conf.FlushesSync() {
			go storage.FlushPeriodic(context.Background())
		}
	}

	controller := usecases.NewBaseController(logger, storage, conf.TemplatePath)
	r := setupServer(logger, controller)
	srv := http.Server{Addr: conf.Addr, Handler: r}
	go func() {
		logger.Infof("Launching the server at %s\n", conf.Addr)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err)
		}
	}()

	gracefulShutdown(&srv, done, logger)
}
