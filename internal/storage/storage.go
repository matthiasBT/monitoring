package storage

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/matthiasBT/monitoring/internal/interfaces"
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

type MemStorage struct {
	MetricsGauge   map[string]float64
	MetricsCounter map[string]int64
	Logger         interfaces.ILogger
}

func (storage *MemStorage) Add(update MetricUpdate) {
	storage.Logger.Infof("Updating metrics with %+v\n", update)
	switch update.Type {
	case TypeGauge:
		storage.Logger.Infof("Old metric value: %f\n", storage.MetricsGauge[update.Name])
		val, _ := strconv.ParseFloat(update.Value, 64)
		storage.MetricsGauge[update.Name] = val
		storage.Logger.Infof("New metric value: %f\n", storage.MetricsGauge[update.Name])
	case TypeCounter:
		storage.Logger.Infof("Old metric value: %d\n", storage.MetricsCounter[update.Name])
		val, _ := strconv.ParseInt(update.Value, 10, 64)
		storage.MetricsCounter[update.Name] += val
		storage.Logger.Infof("New metric value: %d\n", storage.MetricsCounter[update.Name])
	}
}

func (storage *MemStorage) Get(mType string, name string) (string, error) {
	storage.Logger.Infof("Getting metric of type %s named %s\n", mType, name)
	switch mType {
	case TypeGauge:
		if val, ok := storage.MetricsGauge[name]; ok {
			res := strconv.FormatFloat(val, 'f', -1, 64)
			return res, nil
		}
		storage.Logger.Infoln("No such Gauge metric")
		return "", ErrUnknownMetricName
	case TypeCounter:
		if val, ok := storage.MetricsCounter[name]; ok {
			res := fmt.Sprintf("%d", val)
			return res, nil
		}
		storage.Logger.Infoln("No such Counter metric")
		return "", ErrUnknownMetricName
	default:
		storage.Logger.Infoln("Invalid metric type")
		return "", ErrInvalidMetricType
	}
}

func (storage *MemStorage) GetAll() map[string]string {
	res := make(map[string]string, len(storage.MetricsGauge)+len(storage.MetricsCounter))
	for name := range storage.MetricsGauge {
		valStr, _ := storage.Get(TypeGauge, name)
		res[name] = valStr
	}
	for name := range storage.MetricsCounter {
		valStr, _ := storage.Get(TypeCounter, name)
		res[name] = valStr
	}
	return res
}
