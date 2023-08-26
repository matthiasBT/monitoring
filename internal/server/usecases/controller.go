package usecases

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/entities"
)

// TODO: can't split BaseController into 2 parts for moving 1 to adapters and 1 to usecases

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
	r.Post("/update/{type}/{name}/{value}", c.updateMetric)
	r.Get("/value/{type}/{name}", c.getMetric)
	r.Get("/", c.getAllMetrics)
	return r
}

func (c *BaseController) updateMetric(w http.ResponseWriter, r *http.Request) {
	params := extractParams(r, "type", "name", "value")
	UpdateMetric(w, c, params)
}

func (c *BaseController) getMetric(w http.ResponseWriter, r *http.Request) {
	params := extractParams(r, "type", "name")
	GetMetric(w, c, params)
}

func (c *BaseController) getAllMetrics(w http.ResponseWriter, r *http.Request) {
	GetAllMetrics(w, c, "all_metrics.html")
}

func extractParams(r *http.Request, names ...string) map[string]string {
	var params = make(map[string]string)
	for _, name := range names {
		params[name] = chi.URLParam(r, name)
	}
	return params
}
