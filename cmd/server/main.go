package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/matthiasBT/monitoring/internal/adapters"
	"github.com/matthiasBT/monitoring/internal/config"
	"github.com/matthiasBT/monitoring/internal/handlers"
	"github.com/matthiasBT/monitoring/internal/storage"
)

func setupServer() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	return r
}

func main() {
	r := setupServer()
	logger := adapters.SetupLogger()
	conf := config.InitServerConfig(logger)
	controller := handlers.NewBaseController(
		logger,
		&storage.MemStorage{
			MetricsGauge:   make(map[string]float64),
			MetricsCounter: make(map[string]int64),
			Logger:         logger,
		},
		conf.TemplatePath,
	)
	r.Mount("/", controller.Route())
	logger.Fatal(http.ListenAndServe(conf.Addr, r))
}
