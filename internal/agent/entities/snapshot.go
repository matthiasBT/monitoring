package entities

import "github.com/matthiasBT/monitoring/internal/infra/entities"

type Snapshot struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

type SnapshotWrapper struct {
	CurrSnapshot *Snapshot
}

type IReporter interface {
	Report(metrics *entities.Metrics) error
}
