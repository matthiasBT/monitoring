package entities

import (
	"context"

	"github.com/matthiasBT/monitoring/internal/infra/entities"
)

type Storage interface {
	Add(ctx context.Context, update *entities.Metrics) (*entities.Metrics, error)
	Get(ctx context.Context, query *entities.Metrics) (*entities.Metrics, error)
	GetAll(ctx context.Context) (map[string]*entities.Metrics, error)
	AddBatch(ctx context.Context, batch []*entities.Metrics) error
	Snapshot(ctx context.Context) ([]*entities.Metrics, error)
	Init([]*entities.Metrics)
	SetKeeper(keeper Keeper)
	Ping(ctx context.Context) error
}
