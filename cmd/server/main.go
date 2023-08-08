package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/matthiasBT/monitoring/internal/config"
	"github.com/matthiasBT/monitoring/internal/handlers"
	"github.com/matthiasBT/monitoring/internal/storage"
)

func main() {
	conf := config.InitServerConfig()
	r := chi.NewRouter()
	controller := handlers.NewBaseController(&storage.MemStorage{
		MetricsGauge:   make(map[string]float64),
		MetricsCounter: make(map[string]int64),
	})
	r.Mount("/", controller.Route())
	log.Fatal(http.ListenAndServe(conf.Addr, r))
}
