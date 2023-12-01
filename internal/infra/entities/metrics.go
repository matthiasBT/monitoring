// Package entities defines the data structures used for representing
// metrics within the system. It includes the Metrics structure and associated
// logic for validation and formatting.

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
	ErrUnknownMetric     = errors.New("unknown metric")
)

// Metrics represents a single metric data point, including its identifier,
// type, and value. It supports both gauge and counter metric types.
type Metrics struct {
	// ID is the unique identifier of the metric.
	ID string `json:"id"`

	// MType is the type of the metric, such as gauge or counter.
	MType string `json:"type"`

	// Delta is used for counter metrics to represent a change in value.
	Delta *int64 `json:"delta,omitempty"`

	// Value is used for gauge metrics to represent a measurement value.
	Value *float64 `json:"value,omitempty"`
}

// Validate checks the validity of the metric based on its type and the presence
// of the appropriate value (Delta or Value). It ensures that the metric conforms
// to the expected format and contains necessary data.
func (m Metrics) Validate(withValue bool) error {
	if strings.TrimSpace(m.ID) == "" {
		return ErrMissingMetricName
	}
	switch m.MType {
	case TypeGauge:
		if withValue && (m.Value == nil || m.Delta != nil) {
			return ErrInvalidMetricVal
		}
	case TypeCounter:
		if withValue && (m.Delta == nil || m.Value != nil) {
			return ErrInvalidMetricVal
		}
	default:
		return ErrInvalidMetricType
	}
	return nil
}

// ValueAsString returns the metric value as a string, formatted appropriately
// based on the metric type (gauge or counter). For gauge types, it formats
// the value as a float, and for counter types, it formats as an integer.
func (m Metrics) ValueAsString() string {
	var val string
	if m.MType == TypeGauge {
		val = strconv.FormatFloat(*m.Value, 'f', -1, 64)
		if !strings.Contains(val, ".") {
			val += "."
		}
	} else if m.MType == TypeCounter {
		val = fmt.Sprintf("%d", *m.Delta)
	} else {
		return ""
	}
	return val
}
