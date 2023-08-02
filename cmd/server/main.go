package main

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/matthiasBT/monitoring/internal/handlers"
	"github.com/matthiasBT/monitoring/internal/storage"
	"net/http"
)

const addr = ":8080"

var MetricsStorage = storage.MemStorage{
	MetricsGauge:   make(map[string]float64),
	MetricsCounter: make(map[string]int64),
}

func updateMetric(c echo.Context) error {
	err := handlers.UpdateMetric(c, &MetricsStorage)
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

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/update/:type/:name/:value", updateMetric)
	e.Logger.Fatal(e.Start(addr))
}
