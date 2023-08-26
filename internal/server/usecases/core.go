package usecases

import (
	"errors"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/matthiasBT/monitoring/internal/server/entities"
)

func UpdateMetric(w http.ResponseWriter, c *BaseController, params map[string]string) {
	metricUpdate := entities.MetricUpdate{
		Type:  params["type"],
		Name:  params["name"],
		Value: params["value"],
	}
	err := metricUpdate.Validate()
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrInvalidMetricType):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, entities.ErrMissingMetricName):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, entities.ErrInvalidMetricVal):
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		return
	} else {
		c.Stor.Add(metricUpdate)
		w.WriteHeader(http.StatusOK)
		return
	}
}

func GetMetric(w http.ResponseWriter, c *BaseController, params map[string]string) {
	val, err := c.Stor.Get(params["type"], params["name"])
	if err != nil {
		if errors.Is(err, entities.ErrUnknownMetricName) || errors.Is(err, entities.ErrInvalidMetricType) {
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
	data := c.Stor.GetAll()
	path := filepath.Join(c.TemplatePath, templateName)
	tmpl := template.Must(template.ParseFiles(path))
	err := tmpl.Execute(w, data)
	if err != nil {
		c.Logger.Fatal(err)
	}
}
