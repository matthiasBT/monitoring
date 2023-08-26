package entities

import (
	"errors"
	"strconv"
	"strings"
)

const (
	TypeGauge   = "gauge"
	TypeCounter = "counter"
)

type MetricUpdate struct {
	Type  string
	Name  string
	Value string
}

var (
	ErrInvalidMetricType = errors.New("invalid metric type")
	ErrMissingMetricName = errors.New("missing metric name")
	ErrInvalidMetricVal  = errors.New("invalid metric value")
	ErrUnknownMetricName = errors.New("unknown metric name")
)

func (m MetricUpdate) Validate() error {
	if strings.TrimSpace(m.Name) == "" {
		return ErrMissingMetricName
	}
	switch m.Type {
	case TypeGauge:
		if _, err := strconv.ParseFloat(m.Value, 64); err != nil {
			return ErrInvalidMetricVal
		}
	case TypeCounter:
		if _, err := strconv.ParseInt(m.Value, 10, 64); err != nil {
			return ErrInvalidMetricVal
		}
	default:
		return ErrInvalidMetricType
	}
	return nil
}

type Storage interface {
	Add(update MetricUpdate)
	Get(mType string, name string) (string, error)
	GetAll() map[string]string
}
