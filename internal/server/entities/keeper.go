package entities

import "github.com/matthiasBT/monitoring/internal/infra/entities"

type Keeper interface {
	Flush()
	FlushPeriodic()
	Restore() map[string]*entities.Metrics
}
