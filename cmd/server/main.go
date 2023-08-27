package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/matthiasBT/monitoring/internal/infra/compression"
	"github.com/matthiasBT/monitoring/internal/infra/config/server"
	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/adapters"
	"github.com/matthiasBT/monitoring/internal/server/entities"
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

func setupStorage(logger logging.ILogger, events chan<- struct{}) entities.Storage {
	return &adapters.MemStorage{
		Metrics: make(map[string]*common.Metrics),
		Logger:  logger,
		Events:  events,
		Lock:    &sync.Mutex{},
	}
}

func setupFileStorage(
	conf *server.Config,
	logger logging.ILogger,
	storage entities.Storage,
	storageEvents <-chan struct{},
	done chan struct{},
) adapters.FileStorage {
	var tickerChan <-chan time.Time
	if conf.StoresSync() {
		tickerChan = make(chan time.Time) // will never be used
	} else {
		ticker := time.NewTicker(time.Duration(*conf.StoreInterval) * time.Second)
		tickerChan = ticker.C
	}
	return adapters.FileStorage{
		Logger:        logger,
		Storage:       storage,
		Path:          conf.FileStoragePath,
		Done:          done,
		Tick:          tickerChan,
		StorageEvents: storageEvents,
		Lock:          &sync.Mutex{},
		StoreSync:     conf.StoresSync(),
	}
}

func main() {
	logger := logging.SetupLogger()
	conf, err := server.InitConfig()
	if err != nil {
		logger.Fatal(err)
	}

	storageEvents := make(chan struct{})
	storage := setupStorage(logger, storageEvents)

	done := make(chan struct{}, 1)
	fileStorage := setupFileStorage(conf, logger, storage, storageEvents, done)
	if *conf.Restore {
		state := fileStorage.InitStorage()
		storage.Init(state)
	}
	go fileStorage.Dump()

	controller := usecases.NewBaseController(logger, storage, conf.TemplatePath)
	r := setupServer(logger, controller)

	srv := http.Server{Addr: conf.Addr, Handler: r}
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err)
		}
	}()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quitChannel
	logger.Infof("Received signal: %v\n", sig)
	done <- struct{}{}

	time.Sleep(2 * time.Second)

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	}
}
