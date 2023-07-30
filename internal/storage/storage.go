package storage

import (
	"strconv"
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

type Storage interface {
	Add(update MetricUpdate) error
}

type MemStorage struct {
	MetricsGauge   map[string]float64
	MetricsCounter map[string]int64
}

func (storage *MemStorage) Add(update MetricUpdate) {
	switch update.Type {
	case TypeGauge:
		val, _ := strconv.ParseFloat(update.Value, 64)
		storage.MetricsGauge[update.Name] = val
	case TypeCounter:
		val, _ := strconv.ParseInt(update.Value, 10, 64)
		storage.MetricsCounter[update.Name] += val
	}
}
