package entities

import "github.com/matthiasBT/monitoring/internal/infra/entities"

type Storage interface {
	Add(update entities.Metrics) (*entities.Metrics, error)
	Get(query entities.Metrics) (*entities.Metrics, error)
	GetAll() (map[string]*entities.Metrics, error)
	Init(map[string]*entities.Metrics)
}
