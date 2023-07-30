package main

import (
	"errors"
	"github.com/matthiasBT/monitoring/internal/storage"
	"github.com/matthiasBT/monitoring/internal/web"
	"net/http"
)

const addr = ":8080"
const patternUpdate = "/update/"

var MetricsStorage = storage.MemStorage{
	MetricsGauge:   make(map[string]float64),
	MetricsCounter: make(map[string]int64),
}

func updateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	metricUpdate, err := web.ParseMetricUpdate(r.URL.Path, patternUpdate)
	if err == nil {
		MetricsStorage.Add(*metricUpdate)
		w.WriteHeader(http.StatusOK)
		return
	}
	switch {
	case errors.Is(err, web.ErrInvalidMetricType):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, web.ErrMissingMetricName):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, web.ErrInvalidMetricVal):
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc(patternUpdate, updateMetric)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
