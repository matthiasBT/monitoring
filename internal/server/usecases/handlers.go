package usecases

import (
	"errors"
	"net/http"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
)

func (c *BaseController) updateMetric(w http.ResponseWriter, r *http.Request) {
	asJSON := r.Header.Get("Content-Type") == "application/json"
	var metrics *common.Metrics
	if metrics = parseMetric(r, asJSON, true); metrics == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := metrics.Validate(true)
	if err != nil {
		handleInvalidMetric(w, err)
		return
	}

	result, _ := UpdateMetric(c, metrics) // so far, there can't be any errors

	if err := writeMetric(w, asJSON, result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *BaseController) getMetric(w http.ResponseWriter, r *http.Request) {
	asJSON := r.Header.Get("Content-Type") == "application/json"
	var metrics *common.Metrics
	if metrics = parseMetric(r, asJSON, false); metrics == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := metrics.Validate(false)
	if err != nil {
		handleInvalidMetric(w, err)
		return
	}

	result, err := GetMetric(c, metrics)
	if err != nil {
		var status int
		if errors.Is(err, common.ErrUnknownMetricName) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		w.WriteHeader(status)
		w.Write([]byte(err.Error()))
		return
	}

	if err := writeMetric(w, asJSON, result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *BaseController) getAllMetrics(w http.ResponseWriter, r *http.Request) {
	result, _ := GetAllMetrics(c, "all_metrics.html") // so far, there can't be any errors
	w.Write(result.Bytes())
}
