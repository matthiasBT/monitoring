package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/matthiasBT/monitoring/internal/infra/entities"
)

// TODO: move everything web-related to controller

func UpdateMetric(w http.ResponseWriter, c *BaseController, metrics *entities.Metrics) *entities.Metrics {
	err := metrics.Validate()
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
		return nil
	} else {
		updated := c.Stor.Add(*metrics)
		if err := writeMetric(w, true, updated); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return nil
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return updated
	}
}

func GetMetric(w http.ResponseWriter, c *BaseController, asJson bool, metrics *entities.Metrics) {
	metrics, err := c.Stor.Get(metrics.MType, metrics.ID)
	if err != nil {
		if errors.Is(err, entities.ErrUnknownMetricName) || errors.Is(err, entities.ErrInvalidMetricType) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))
		return
	} else {
		if asJson {
			w.Header().Set("Content-Type", "application/json")
		}
		if err := writeMetric(w, asJson, metrics); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
}

// maybe this handler is no longer needed

func GetAllMetrics(w http.ResponseWriter, c *BaseController, templateName string) {
	metrics := c.Stor.GetAll()
	var data = make(map[string]string, len(metrics))
	for _, m := range metrics {
		var val string
		if m.MType == entities.TypeGauge {
			val = strconv.FormatFloat(*m.Value, 'f', -1, 64)
		} else {
			val = fmt.Sprintf("%d", *m.Delta)
		}
		data[m.ID] = val
	}
	path := filepath.Join(c.TemplatePath, templateName)
	tmpl := template.Must(template.ParseFiles(path))
	err := tmpl.Execute(w, data)
	if err != nil {
		c.Logger.Fatal(err)
	}
}

func writeMetric(w http.ResponseWriter, asJson bool, metrics *entities.Metrics) error {
	var body []byte
	if asJson {
		val, err := json.Marshal(metrics)
		if err != nil {
			return err
		}
		body = val
	} else {
		body = []byte(valToStr(metrics))
	}
	w.Write(body)
	return nil
}

func valToStr(metrics *entities.Metrics) string {
	var val string
	if metrics.MType == entities.TypeGauge {
		val = strconv.FormatFloat(*metrics.Value, 'f', -1, 64)
	} else {
		val = fmt.Sprintf("%d", *metrics.Delta)
	}
	return val
}
