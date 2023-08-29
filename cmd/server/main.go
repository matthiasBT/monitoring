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
		log.Fatalf("Server shutdown failed: %v\n", err)
	}
}

func main() {
	logger := logging.SetupLogger()
	conf, err := server.InitConfig()
	if err != nil {
		logger.Fatal(err)
	}

	storageEvents := make(chan struct{})
	storage := adapters.NewMemStorage(logger, storageEvents)

	done := make(chan struct{}, 1)
	fileStorage := adapters.NewFileStorage(conf, logger, storage, storageEvents, done)
	if *conf.Restore {
		state := fileStorage.InitStorage()
		storage.Init(state)
	}
	go fileStorage.Flush()

	controller := usecases.NewBaseController(logger, storage, conf.TemplatePath)
	r := setupServer(logger, controller)
	srv := http.Server{Addr: conf.Addr, Handler: r}
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err)
		}
	}()

	gracefulShutdown(&srv, done, logger)
}
