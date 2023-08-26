package usecases

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	common "github.com/matthiasBT/monitoring/internal/infra/entities"
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
	r.Post("/update/", c.updateMetric)
	r.Post("/value/", c.getMetric)
	r.Post("/update/{type}/{name}/{value}", c.updateMetric)
	r.Get("/value/{type}/{name}", c.getMetric)
	r.Get("/", c.getAllMetrics)
	return r
}

func (c *BaseController) updateMetric(w http.ResponseWriter, r *http.Request) {
	asJson := r.Header.Get("Content-Type") == "application/json"
	if metrics := parseMetric(w, r, asJson, false); metrics != nil {
		UpdateMetric(w, c, metrics)
	}
}

func (c *BaseController) getMetric(w http.ResponseWriter, r *http.Request) {
	asJson := r.Header.Get("Content-Type") == "application/json"
	if metrics := parseMetric(w, r, asJson, true); metrics != nil {
		GetMetric(w, c, asJson, metrics)
	}
}

func (c *BaseController) getAllMetrics(w http.ResponseWriter, r *http.Request) {
	GetAllMetrics(w, c, "all_metrics.html")
}

func parseMetric(w http.ResponseWriter, r *http.Request, asjson bool, withoutValue bool) *common.Metrics {
	var metrics common.Metrics
	if asjson {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
		if err := json.Unmarshal(body, &metrics); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		metrics.ID = chi.URLParam(r, "name")
		metrics.MType = chi.URLParam(r, "type")
		if withoutValue {
			return &metrics
		}
		val := chi.URLParam(r, "value")
		switch metrics.MType {
		case common.TypeGauge:
			val, err := strconv.ParseFloat(val, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			metrics.Value = &val
		case common.TypeCounter:
			val, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			metrics.Delta = &val
		default:
			w.WriteHeader(http.StatusBadRequest)
			return nil
		}
	}
	return &metrics
}
