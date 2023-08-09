package handlers

import (
	"errors"
	"html/template"
	"net/http"
	"path/filepath"

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
	switch {
	case errors.Is(err, storage.ErrInvalidMetricType):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, storage.ErrMissingMetricName):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, storage.ErrInvalidMetricVal):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(err.Error()))
}

func GetMetric(w http.ResponseWriter, c *BaseController, params map[string]string) {
	val, err := c.stor.Get(params["type"], params["name"])
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(val))
		return
	}
	if errors.Is(err, storage.ErrUnknownMetricName) || errors.Is(err, storage.ErrInvalidMetricType) {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(err.Error()))
}

func GetAllMetrics(w http.ResponseWriter, c *BaseController, templateName string) {
	data := c.stor.GetAll()
	path := filepath.Join(c.templatePath, templateName)
	tmpl := template.Must(template.ParseFiles(path))
	tmpl.Execute(w, data) // todo: handle error
}
