// Package usecases provides functions that encapsulate the core business logic
// for handling metrics within the monitoring application. It includes
// functions for updating and retrieving metrics, handling batch updates,
// and rendering metrics for presentation.
package usecases

import (
	"bytes"
	"context"
	"html/template"
	"path/filepath"

	"github.com/matthiasBT/monitoring/internal/infra/entities"
)

// UpdateMetric updates a single metric in the storage using the provided BaseController.
// It returns the updated metric and an error, if any.
func UpdateMetric(ctx context.Context, c *BaseController, metrics *entities.Metrics) (*entities.Metrics, error) {
	result, err := c.Stor.Add(ctx, metrics)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetMetric retrieves a single metric from the storage based on the provided query criteria.
// It uses the BaseController and returns the found metric and an error, if any.
func GetMetric(ctx context.Context, c *BaseController, metrics *entities.Metrics) (*entities.Metrics, error) {
	return c.Stor.Get(ctx, metrics)
}

// GetAllMetrics retrieves all metrics from the storage and renders them using a specified HTML template.
// The rendered data is returned as a bytes.Buffer, along with an error, if any.
func GetAllMetrics(ctx context.Context, c *BaseController, templateName string) (*bytes.Buffer, error) {
	metrics, err := c.Stor.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	data := prepareTemplateData(metrics)
	initialBytes := make([]byte, 0, len(data))
	var result bytes.Buffer
	result.Write(initialBytes)
	path := filepath.Join(c.TemplatePath, templateName)
	tmpl := template.Must(template.ParseFiles(path))
	if err := tmpl.Execute(&result, data); err != nil {
		return nil, err
	}
	return &result, nil
}

// MassUpdate updates a batch of metrics in the storage using the provided BaseController.
// It returns an error, if any occurs during the operation.
func MassUpdate(ctx context.Context, c *BaseController, batch []*entities.Metrics) error {
	return c.Stor.AddBatch(ctx, batch)
}

// prepareTemplateData prepares a map of metrics data for rendering in an HTML template.
// It converts the metrics data into a suitable format for templating.
func prepareTemplateData(metrics map[string]*entities.Metrics) map[string]string {
	var data = make(map[string]string, len(metrics))
	for _, m := range metrics {
		data[m.ID] = m.ValueAsString()
	}
	return data
}
