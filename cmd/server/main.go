package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/matthiasBT/monitoring/internal/infra/config/server"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/adapters"
	"github.com/matthiasBT/monitoring/internal/server/usecases"
)

func setupServer() *chi.Mux {
	r := chi.NewRouter()
	return r
}

func main() {
	r := setupServer()
	logger := logging.SetupLogger()
	conf, err := server.InitServerConfig()
	if err != nil {
		logger.Fatal(err)
	}
	controller := usecases.NewBaseController(
		logger,
		&adapters.MemStorage{
			MetricsGauge:   make(map[string]float64),
			MetricsCounter: make(map[string]int64),
			Logger:         logger,
		},
		conf.TemplatePath,
	)
	r.Mount("/", controller.Route())
	logger.Fatal(http.ListenAndServe(conf.Addr, r))
}
