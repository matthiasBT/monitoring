package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	conf := config.InitServerConfig()
	r := setupServer()
	controller := handlers.NewBaseController(
		&storage.MemStorage{
			MetricsGauge:   make(map[string]float64),
			MetricsCounter: make(map[string]int64),
		},
		conf.TemplatePath,
	)
	r.Mount("/", controller.Route())
	log.Fatal(http.ListenAndServe(conf.Addr, r))
}
