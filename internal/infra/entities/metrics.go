package entities

import (
	"errors"
	"strings"
)

const (
	TypeGauge   = "gauge"
	TypeCounter = "counter"
)

var (
	ErrInvalidMetricType = errors.New("invalid metric type")
	ErrMissingMetricName = errors.New("missing metric name")
	ErrInvalidMetricVal  = errors.New("invalid metric value")
	ErrUnknownMetricName = errors.New("unknown metric name")
)

type Metrics struct {
	ID    string   `easyjson:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m Metrics) Validate() error {
	if strings.TrimSpace(m.ID) == "" {
		return ErrMissingMetricName
	}
	switch m.MType {
	case TypeGauge:
		if m.Value == nil {
			return ErrInvalidMetricVal
		}
	case TypeCounter:
		if m.Delta == nil {
			return ErrInvalidMetricVal
		}
	default:
		return ErrInvalidMetricType
	}
	return nil
}
