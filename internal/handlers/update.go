package handlers

import (
	"errors"
	"fmt"
	"github.com/matthiasBT/monitoring/internal/storage"
	"net/http"
)

func UpdateMetric(w http.ResponseWriter, r *http.Request, patternUpdate string, stor *storage.MemStorage) {
	fmt.Printf("Request: %v %v\n", r.Method, r.URL)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	metricUpdate, err := parseMetricUpdate(r.URL.Path, patternUpdate)
	if err == nil {
		stor.Add(*metricUpdate)
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
