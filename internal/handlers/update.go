package handlers

import (
	"errors"
	"github.com/matthiasBT/monitoring/internal/storage"
	"net/http"
)

func UpdateMetric(w http.ResponseWriter, r *http.Request, patternUpdate string, storage storage.MemStorage) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	metricUpdate, err := parseMetricUpdate(r.URL.Path, patternUpdate)
	if err == nil {
		storage.Add(*metricUpdate)
		w.WriteHeader(http.StatusOK)
		return
	}
	switch {
	case errors.Is(err, ErrInvalidMetricType):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrMissingMetricName):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInvalidMetricVal):
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
