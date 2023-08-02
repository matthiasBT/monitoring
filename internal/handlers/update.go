package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/matthiasBT/monitoring/internal/storage"
	"net/http"
)

func UpdateMetric(c echo.Context, stor *storage.MemStorage) error {
	metricUpdate := storage.MetricUpdate{
		Type:  c.Param("type"),
		Name:  c.Param("name"),
		Value: c.Param("value"),
	}
	if err := metricUpdate.Validate(); err != nil {
		return err
	}
	stor.Add(metricUpdate)
	c.Response().WriteHeader(http.StatusOK)
	return nil
}
