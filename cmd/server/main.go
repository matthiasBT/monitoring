package main

import (
	"github.com/matthiasBT/monitoring/internal/handlers"
	"github.com/matthiasBT/monitoring/internal/storage"
	"net/http"
)

const addr = ":8080"
const patternUpdate = "/update/"

var MetricsStorage = storage.MemStorage{
	MetricsGauge:   make(map[string]float64),
	MetricsCounter: make(map[string]int64),
}

func updateMetric(w http.ResponseWriter, r *http.Request) {
	handlers.UpdateMetric(w, r, patternUpdate, MetricsStorage)
}

func main() {
	http.HandleFunc(patternUpdate, updateMetric)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
