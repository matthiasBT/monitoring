package handlers

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/matthiasBT/monitoring/internal/storage"
)

func UpdateMetric(w http.ResponseWriter, c *BaseController, params map[string]string) {
	metricUpdate := storage.MetricUpdate{
		Type:  params["type"],
		Name:  params["name"],
		Value: params["value"],
	}
	err := metricUpdate.Validate()
	if err == nil {
		c.stor.Add(metricUpdate)
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Write([]byte(err.Error()))
	switch {
	case errors.Is(err, storage.ErrInvalidMetricType):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, storage.ErrMissingMetricName):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, storage.ErrInvalidMetricVal):
		w.WriteHeader(http.StatusBadRequest)
	}
}

func GetMetric(w http.ResponseWriter, c *BaseController, params map[string]string) {
	val, err := c.stor.Get(params["type"], params["name"])
	if err == nil {
		w.Write([]byte(val))
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Write([]byte(err.Error()))
	if errors.Is(err, storage.ErrUnknownMetricName) || errors.Is(err, storage.ErrInvalidMetricType) {
		w.WriteHeader(http.StatusNotFound)
	}
}

func GetAllMetrics(w http.ResponseWriter, c *BaseController) {
	data := c.stor.GetAll()
	tmpl := template.Must(template.ParseFiles("web/template/all_metrics.html"))
	tmpl.Execute(w, data) // todo: handle error
}
