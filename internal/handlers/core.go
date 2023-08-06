package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/matthiasBT/monitoring/internal/storage"
	"net/http"
)

func UpdateMetric(c echo.Context, stor storage.Storage) error {
	metricUpdate := storage.MetricUpdate{
		Type:  c.Param("type"),
		Name:  c.Param("name"),
		Value: c.Param("value"),
	}
	err := metricUpdate.Validate()
	if err == nil {
		stor.Add(metricUpdate)
		c.Response().WriteHeader(http.StatusOK)
		return nil
	}
	switch {
	case errors.Is(err, storage.ErrInvalidMetricType):
		c.String(http.StatusBadRequest, err.Error())
	case errors.Is(err, storage.ErrMissingMetricName):
		c.String(http.StatusNotFound, err.Error())
	case errors.Is(err, storage.ErrInvalidMetricVal):
		c.String(http.StatusBadRequest, err.Error())
	default:
		return err
	}
	return nil
}

func GetMetric(c echo.Context, stor storage.Storage) error {
	mType := c.Param("type")
	name := c.Param("name")
	val, err := stor.Get(mType, name)
	if err == nil {
		c.String(http.StatusOK, val)
		return nil
	}
	switch {
	case errors.Is(err, storage.ErrUnknownMetricName) || errors.Is(err, storage.ErrInvalidMetricType):
		c.String(http.StatusNotFound, err.Error())
	default:
		return err
	}
	return nil
}

func GetAllMetrics(c echo.Context, stor storage.Storage) error {
	res := stor.GetAll()
	return c.Render(http.StatusOK, "all_metrics", res)
}
