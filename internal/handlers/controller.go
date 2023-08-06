package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/matthiasBT/monitoring/internal/storage"
)

type BaseController struct {
	e    *echo.Echo
	stor storage.Storage
}

func NewBaseController(e *echo.Echo, stor storage.Storage) *BaseController {
	return &BaseController{e: e, stor: stor}
}

func (c *BaseController) Route(prefix string) *echo.Group {
	g := c.e.Group(prefix)
	g.POST("/update/:type/:name/:value", c.updateMetric)
	g.GET("/value/:type/:name", c.getMetric)
	g.GET("/", c.getAllMetrics)
	return g
}

func (c *BaseController) updateMetric(ctx echo.Context) error {
	return UpdateMetric(ctx, c.stor)
}

func (c *BaseController) getMetric(ctx echo.Context) error {
	return GetMetric(ctx, c.stor)
}

func (c *BaseController) getAllMetrics(ctx echo.Context) error {
	return GetAllMetrics(ctx, c.stor)
}
