package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/matthiasBT/monitoring/internal/handlers"
	"github.com/matthiasBT/monitoring/internal/storage"
	"html/template"
	"io"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

const addr = ":8080"

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
	e := echo.New()
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("web/template/*.html")),
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/update/:type/:name/:value", updateMetric)
	e.GET("/value/:type/:name", getMetric)
	e.GET("/", getAllMetrics)
	e.Logger.Fatal(e.Start(addr))
}
