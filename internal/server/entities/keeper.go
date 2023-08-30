package entities

import "github.com/matthiasBT/monitoring/internal/infra/entities"

type Keeper interface {
	Flush([]*entities.Metrics) error
	Restore() []*entities.Metrics
}
