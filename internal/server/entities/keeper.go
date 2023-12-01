// Package entities defines interfaces and types for abstracting
// storage operations in the monitoring application.
package entities

import (
	"context"

	"github.com/matthiasBT/monitoring/internal/infra/entities"
)

// Keeper is an interface that abstracts operations for storing and retrieving
// metrics data. It defines methods for flushing data to storage, restoring data
// from storage, checking the storage status, and performing shutdown operations.
type Keeper interface {
	// Flush writes a slice of Metrics to the storage. It may involve
	// complex operations like database transactions or file operations,
	// depending on the implementation.
	Flush(context.Context, []*entities.Metrics) error

	// Restore retrieves all stored metrics from the storage. The method
	// is expected to return a slice of Metrics, potentially involving
	// database queries or file read operations.
	Restore() []*entities.Metrics

	// Ping checks the liveness or connectivity of the storage system.
	// It's particularly useful in scenarios like database connections
	// where the storage might be remote or require a live network connection.
	Ping(ctx context.Context) error

	// Shutdown gracefully terminates any connections or operations related
	// to the storage. It's essential for releasing resources or closing
	// connections in a managed way.
	Shutdown()
}
