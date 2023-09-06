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

	"github.com/matthiasBT/monitoring/cmd/server/periodic"
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

func main() {
	logger := logging.SetupLogger()
	conf, err := server.InitConfig()
	if err != nil {
		logger.Fatal(err)
	}

	var DBManager *adapters.DBManager
	if conf.DatabaseDSN != "" {
		DBManager, err = adapters.NewDBManager(conf.DatabaseDSN, logger)
		if err != nil {
			logger.Errorf("Failed to init database connection: %s\n", err.Error())
			panic(err)
		}
		defer DBManager.Shutdown()
	}

	var storage entities.Storage
	done := make(chan struct{}, 1)
	if DBManager != nil {
		storage = adapters.NewDBStorage(DBManager.DB, logger, nil)
	} else {
		storage = adapters.NewMemStorage(logger, nil)
		if conf.Flushes() {
			fileKeeper := adapters.NewFileKeeper(conf, logger, done)
			flusher := periodic.NewFlusher(conf, logger, storage, fileKeeper, done)
			if *conf.Restore {
				state := fileKeeper.Restore()
				storage.Init(state)
			}
			if conf.FlushesSync() {
				storage.SetKeeper(fileKeeper)
			} else {
				go flusher.Flush(context.Background())
			}
		}
	}

	controller := usecases.NewBaseController(logger, storage, DBManager, conf.TemplatePath)
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
