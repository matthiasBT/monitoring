package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/matthiasBT/monitoring/internal/infra/config/server"
	"github.com/matthiasBT/monitoring/internal/infra/entities"
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
	conf, err := server.InitConfig()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Server config: %v\n", *conf)
	controller := usecases.NewBaseController(
		logger,
		&adapters.MemStorage{
			Metrics: make(map[string]*entities.Metrics),
			Logger:  logger,
		},
		conf.TemplatePath,
	)
	r.Use(logging.Middleware(logger))
	r.Mount("/", controller.Route())
	logger.Fatal(http.ListenAndServe(conf.Addr, r))
}
