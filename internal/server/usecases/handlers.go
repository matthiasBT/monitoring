// Package usecases provides methods for handling HTTP requests related to
// metrics management in the monitoring application. It includes methods for
// updating, retrieving, and batch processing metrics, as well as health checking.
package usecases

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
)

// updateMetric handles the HTTP request for updating a metric.
// It supports both JSON and form data, validates the input, and writes
// the updated metric back to the response.
func (c *BaseController) updateMetric(w http.ResponseWriter, r *http.Request) {
	asJSON := r.Header.Get("Content-Type") == "application/json"
	var metrics *common.Metrics
	if metrics = parseMetric(r, asJSON, true); metrics == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse metric"))
		return
	}

	if err := metrics.Validate(true); err != nil {
		handleInvalidMetric(w, err)
		return
	}

	result, err := UpdateMetric(r.Context(), c, metrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := writeMetric(w, asJSON, result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

// getMetric handles the HTTP request for retrieving a specific metric.
// It supports both JSON and form data, validates the query, and writes
// the metric data back to the response.
func (c *BaseController) getMetric(w http.ResponseWriter, r *http.Request) {
	asJSON := r.Header.Get("Content-Type") == "application/json"
	var metrics *common.Metrics
	if metrics = parseMetric(r, asJSON, false); metrics == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse metric"))
		return
	}

	err := metrics.Validate(false)
	if err != nil {
		handleInvalidMetric(w, err)
		return
	}

	result, err := GetMetric(r.Context(), c, metrics)
	if err != nil {
		var status int
		if errors.Is(err, common.ErrUnknownMetric) {
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
		w.Write([]byte(err.Error()))
	}
}

// getAllMetrics handles the HTTP request for retrieving all metrics.
// It renders the metrics in an HTML template and sends the result back to the response.
func (c *BaseController) getAllMetrics(w http.ResponseWriter, r *http.Request) {
	result, err := GetAllMetrics(r.Context(), c, "all_metrics.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(result.Bytes())
}

// massUpdate handles the HTTP request for updating a batch of metrics.
// It only accepts JSON data, validates the input, and sends an appropriate response.
func (c *BaseController) massUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Supply data as JSON"))
		return
	}

	var batch []*common.Metrics
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err := json.Unmarshal(body, &batch); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	for _, metrics := range batch {
		if err := metrics.Validate(true); err != nil {
			handleInvalidMetric(w, err)
			return
		}
	}

	if err := MassUpdate(r.Context(), c, batch); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

// ping handles the HTTP request for checking the storage connectivity or liveliness.
// It uses the Ping method of the storage and sends an appropriate response.
func (c *BaseController) ping(w http.ResponseWriter, r *http.Request) {
	if err := c.Stor.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
}
