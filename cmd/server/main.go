package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/matthiasBT/monitoring/internal/config"
	"github.com/matthiasBT/monitoring/internal/handlers"
	"github.com/matthiasBT/monitoring/internal/storage"
)

const templatePath = "web/template/*.html"

func setupServer() *echo.Echo {
	e := echo.New()
	e.Renderer = handlers.GetRenderer(templatePath)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	return e
}

func main() {
	conf := config.InitServerConfig()
	e := setupServer()
	c := handlers.NewBaseController(e, &storage.MemStorage{
		MetricsGauge:   make(map[string]float64),
		MetricsCounter: make(map[string]int64),
	})
	c.Route("")
	e.Logger.Fatal(e.Start(conf.Addr))
}
