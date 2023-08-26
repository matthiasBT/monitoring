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
	if err != nil {
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
		return
	} else {
		c.stor.Add(metricUpdate)
		w.WriteHeader(http.StatusOK)
		return
	}
}

func GetMetric(w http.ResponseWriter, c *BaseController, params map[string]string) {
	val, err := c.stor.Get(params["type"], params["name"])
	if err != nil {
		if errors.Is(err, storage.ErrUnknownMetricName) || errors.Is(err, storage.ErrInvalidMetricType) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(val))
		return
	}
}

func GetAllMetrics(w http.ResponseWriter, c *BaseController, templateName string) {
	data := c.stor.GetAll()
	path := filepath.Join(c.templatePath, templateName)
	tmpl := template.Must(template.ParseFiles(path))
	err := tmpl.Execute(w, data)
	if err != nil {
		c.logger.Fatal(err)
	}
}
