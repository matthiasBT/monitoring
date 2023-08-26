package entities

import "github.com/matthiasBT/monitoring/internal/infra/entities"

type Storage interface {
	Add(update entities.Metrics) *entities.Metrics
	Get(mType string, name string) (*entities.Metrics, error)
	GetAll() map[string]*entities.Metrics
}
