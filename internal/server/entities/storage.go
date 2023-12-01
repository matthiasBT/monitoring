// Package entities defines interfaces and types for abstracting
// storage operations in the monitoring application. This package
// specifically focuses on the Storage interface for managing metrics data.
package entities

import (
	"context"

	"github.com/matthiasBT/monitoring/internal/infra/entities"
)

// Storage is an interface that abstracts operations for storing, retrieving,
// and managing metrics data. It defines a set of methods for handling metrics
// in various ways, including adding individual or batches of metrics,
// retrieving them, and maintaining the overall state of the storage.
type Storage interface {
	// Add inserts or updates a single metric in the storage.
	// Returns the updated metric and an error, if any.
	Add(ctx context.Context, update *entities.Metrics) (*entities.Metrics, error)

	// Get retrieves a specific metric based on the provided query criteria.
	// Returns the found metric and an error, if any.
	Get(ctx context.Context, query *entities.Metrics) (*entities.Metrics, error)

	// GetAll returns all the metrics currently stored.
	// Returns a map of metrics and an error, if any.
	GetAll(ctx context.Context) (map[string]*entities.Metrics, error)

	// AddBatch inserts or updates a batch of metrics in the storage.
	// Returns an error, if any occurs during the operation.
	AddBatch(ctx context.Context, batch []*entities.Metrics) error

	// Snapshot creates and returns a snapshot of the current metrics in storage.
	// This is typically used for backup or synchronization purposes.
	Snapshot(ctx context.Context) ([]*entities.Metrics, error)

	// Init initializes the storage with a given set of metrics.
	// This could be used for setting up the storage with default or initial data.
	Init([]*entities.Metrics)

	// Ping checks the liveness or connectivity of the storage system,
	// ensuring that it is ready for operations.
	Ping(ctx context.Context) error

	// FlushPeriodic handles the periodic flushing of data from the storage
	// to some persistent or external storage system. This is often used
	// to ensure data durability and consistency.
	FlushPeriodic(ctx context.Context)
}
