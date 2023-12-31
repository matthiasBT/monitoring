// Package usecases provides utility functions for parsing and writing metrics data
// in HTTP requests and responses. These functions support both JSON and URL-encoded data.
package usecases

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	common "github.com/matthiasBT/monitoring/internal/infra/entities"
)

// parseMetric parses a metric from an HTTP request. It supports JSON and URL-encoded data.
// For URL-encoded data, it can parse with or without the metric value.
func parseMetric(r *http.Request, asJSON bool, withValue bool) *common.Metrics {
	var metrics common.Metrics
	if asJSON {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil
		}
		if err := json.Unmarshal(body, &metrics); err != nil {
			return nil
		}
	} else {
		metrics.ID = chi.URLParam(r, "name")
		metrics.MType = chi.URLParam(r, "type")
		if !withValue {
			return &metrics
		}
		val := chi.URLParam(r, "value")
		switch metrics.MType {
		case common.TypeGauge:
			val, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil
			}
			metrics.Value = &val
		case common.TypeCounter:
			val, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil
			}
			metrics.Delta = &val
		}
	}
	return &metrics
}

// writeMetric writes a metric to an HTTP response. It supports both JSON and plain text formats.
func writeMetric(w http.ResponseWriter, asJSON bool, metrics *common.Metrics) error {
	var body []byte
	if asJSON {
		val, err := json.Marshal(metrics)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		body = val
	} else {
		body = []byte(metrics.ValueAsString())
	}
	w.Write(body)
	return nil
}

// handleInvalidMetric handles errors related to invalid metrics and writes appropriate
// HTTP status codes and messages to the response. It covers various error cases like
// invalid metric type, missing metric name, or invalid metric value.
func handleInvalidMetric(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, common.ErrInvalidMetricType):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, common.ErrMissingMetricName):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, common.ErrInvalidMetricVal):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(err.Error()))
}
