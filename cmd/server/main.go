package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/matthiasBT/monitoring/internal/config"
	"github.com/matthiasBT/monitoring/internal/handlers"
	"github.com/matthiasBT/monitoring/internal/storage"
)

var metricsStorage = storage.MemStorage{
	MetricsGauge:   make(map[string]float64),
	MetricsCounter: make(map[string]int64),
}

func updateMetric(c echo.Context) error {
	return handlers.UpdateMetric(c, &metricsStorage)
}

func getMetric(c echo.Context) error {
	return handlers.GetMetric(c, &metricsStorage)
}

func getAllMetrics(c echo.Context) error {
	return handlers.GetAllMetrics(c, &metricsStorage)
}

func main() {
	conf := config.InitServerConfig()
	e := echo.New()
	e.Renderer = handlers.GetRenderer("web/template/*.html")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/update/:type/:name/:value", updateMetric)
	e.GET("/value/:type/:name", getMetric)
	e.GET("/", getAllMetrics)
	e.Logger.Fatal(e.Start(conf.Addr))
}
