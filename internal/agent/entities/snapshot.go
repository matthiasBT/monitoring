// Package entities defines data structures and interfaces used for representing
// and reporting monitoring data in the monitoring system.
package entities

import "github.com/matthiasBT/monitoring/internal/infra/entities"

// Snapshot represents a snapshot of monitoring data at a specific point in time.
// It includes gauges and counters, where gauges represent measurements at a particular
// moment and counters represent cumulative values over time.
type Snapshot struct {
	// Gauges is a map where keys are gauge names and values are the gauge measurements.
	Gauges map[string]float64

	// Counters is a map where keys are counter names and values are the accumulated counter values.
	Counters map[string]int64
}

// SnapshotWrapper encapsulates a Snapshot. It is used for passing snapshots
// around in the system, potentially adding more metadata or context in the future.
type SnapshotWrapper struct {
	// CurrSnapshot is a pointer to the current Snapshot being encapsulated.
	CurrSnapshot *Snapshot
}

// IReporter is an interface defining methods for reporting monitoring metrics.
// Implementations of this interface are responsible for handling the actual
// reporting of metric data.
type IReporter interface {
	// ReportBatch sends a batch of metrics. This method is useful for reporting
	// multiple metrics at once, potentially reducing overhead or network calls.
	ReportBatch(batch []*entities.Metrics) error
}
