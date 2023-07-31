package storage

import (
	"fmt"
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
