package entities

import (
	"context"

	"github.com/matthiasBT/monitoring/internal/infra/entities"
)

type Keeper interface {
	Flush(context.Context, []*entities.Metrics) error
	Restore() []*entities.Metrics
	Ping(ctx context.Context) error
	Shutdown()
}
