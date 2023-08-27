package entities

import (
	"errors"
	"fmt"
	"strconv"
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

func (m Metrics) Validate(withValue bool) error {
	if strings.TrimSpace(m.ID) == "" {
		return ErrMissingMetricName
	}
	switch m.MType {
	case TypeGauge:
		if withValue && m.Value == nil {
			return ErrInvalidMetricVal
		}
	case TypeCounter:
		if withValue && m.Delta == nil {
			return ErrInvalidMetricVal
		}
	default:
		return ErrInvalidMetricType
	}
	return nil
}

func (m Metrics) ValueAsString() string {
	var val string
	if m.MType == TypeGauge {
		val = strconv.FormatFloat(*m.Value, 'f', -1, 64)
	} else if m.MType == TypeCounter {
		val = fmt.Sprintf("%d", *m.Delta)
	} else {
		return ""
	}
	return val
}
