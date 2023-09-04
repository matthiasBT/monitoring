package usecases

import (
	"bytes"
	"context"
	"html/template"
	"path/filepath"

	"github.com/matthiasBT/monitoring/internal/infra/entities"
)

func UpdateMetric(ctx context.Context, c *BaseController, metrics *entities.Metrics) (*entities.Metrics, error) {
	result, err := c.Stor.Add(ctx, *metrics)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetMetric(ctx context.Context, c *BaseController, metrics *entities.Metrics) (*entities.Metrics, error) {
	result, err := c.Stor.Get(ctx, *metrics)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetAllMetrics(ctx context.Context, c *BaseController, templateName string) (*bytes.Buffer, error) {
	metrics, err := c.Stor.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var result bytes.Buffer
	data := prepareTemplateData(metrics)
	path := filepath.Join(c.TemplatePath, templateName)
	tmpl := template.Must(template.ParseFiles(path))
	if err := tmpl.Execute(&result, data); err != nil {
		return nil, err
	}
	return &result, nil
}

func prepareTemplateData(metrics map[string]*entities.Metrics) map[string]string {
	var data = make(map[string]string, len(metrics))
	for _, m := range metrics {
		data[m.ID] = m.ValueAsString()
	}
	return data
}
