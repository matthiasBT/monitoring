package usecases

import (
	"github.com/go-chi/chi/v5"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/entities"
)

type BaseController struct {
	Logger       logging.ILogger
	Stor         entities.Storage
	TemplatePath string
}

func NewBaseController(logger logging.ILogger, stor entities.Storage, templatePath string) *BaseController {
	return &BaseController{logger, stor, templatePath}
}

func (c *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/update/", c.updateMetric)
	r.Post("/value/", c.getMetric)
	r.Post("/update/{type}/{name}/{value}", c.updateMetric)
	r.Get("/value/{type}/{name}", c.getMetric)
	r.Get("/", c.getAllMetrics)
	return r
}
