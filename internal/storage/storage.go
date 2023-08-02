package storage

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
	Add(update MetricUpdate) error
	Get(mType string, name string) (*string, error)
}

type MemStorage struct {
	MetricsGauge   map[string]float64
	MetricsCounter map[string]int64
}

func (storage *MemStorage) Add(update MetricUpdate) {
	fmt.Printf("Updating metrics with %+v\n", update)
	switch update.Type {
	case TypeGauge:
		fmt.Printf("Old metric value: %f\n", storage.MetricsGauge[update.Name])
		val, _ := strconv.ParseFloat(update.Value, 64)
		storage.MetricsGauge[update.Name] = val
		fmt.Printf("New metric value: %f\n", storage.MetricsGauge[update.Name])
	case TypeCounter:
		fmt.Printf("Old metric value: %d\n", storage.MetricsCounter[update.Name])
		val, _ := strconv.ParseInt(update.Value, 10, 64)
		storage.MetricsCounter[update.Name] += val
		fmt.Printf("New metric value: %d\n", storage.MetricsCounter[update.Name])
	}
}

func (storage *MemStorage) Get(mType string, name string) (*string, error) {
	fmt.Printf("Getting metric of type %s named %s\n", mType, name)
	switch mType {
	case TypeGauge:
		if val, ok := storage.MetricsGauge[name]; ok {
			res := strconv.FormatFloat(val, 'f', -1, 64)
			return &res, nil
		}
		fmt.Printf("No such Gauge metric")
		return nil, ErrUnknownMetricName
	case TypeCounter:
		if val, ok := storage.MetricsCounter[name]; ok {
			res := fmt.Sprintf("%d", val)
			return &res, nil
		}
		fmt.Printf("No such Counter metric")
		return nil, ErrUnknownMetricName
	default:
		fmt.Printf("Invalid metric type")
		return nil, ErrInvalidMetricType
	}
}
