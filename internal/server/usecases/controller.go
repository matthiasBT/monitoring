// Package usecases provides the controller layer for handling HTTP requests.
// It includes the BaseController struct which sets up routing for various
// endpoints related to metrics operations.
package usecases

import (
	"github.com/go-chi/chi/v5"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/entities"
)

// BaseController is a struct that holds a logger, storage interface, and a path to HTML templates.
// It is responsible for handling HTTP requests and directing them to appropriate handlers.
type BaseController struct {
	Logger       logging.ILogger  // Logger for logging activities
	Stor         entities.Storage // Storage interface for managing metrics data
	TemplatePath string           // Path to HTML templates
}

// NewBaseController creates and returns a new instance of BaseController.
// It initializes the controller with a logger, storage interface, and template path.
func NewBaseController(logger logging.ILogger, stor entities.Storage, templatePath string) *BaseController {
	return &BaseController{
		Logger:       logger,
		Stor:         stor,
		TemplatePath: templatePath,
	}
}

// Route sets up the HTTP routes for the BaseController. It defines endpoints
// for operations like pinging the server, updating metrics, retrieving metrics,
// batch updating metrics, and retrieving all metrics.
func (c *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/ping", c.ping)
	r.Post("/update/", c.updateMetric)
	r.Post("/value/", c.getMetric)
	r.Post("/update/{type}/{name}/{value}", c.updateMetric)
	r.Get("/value/{type}/{name}", c.getMetric)
	r.Post("/updates/", c.massUpdate)
	r.Get("/", c.getAllMetrics)
	return r
}
