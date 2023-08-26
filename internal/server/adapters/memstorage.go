package adapters

import (
	"github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
)

// TODO: store everything in a single map with Metrics entity

type MemStorage struct {
	Metrics map[string]*entities.Metrics
	Logger  logging.ILogger
}

// TODO: add error

func (storage *MemStorage) Add(update entities.Metrics) *entities.Metrics {
	storage.Logger.Infof("Updating a metric %s %s\n", update.ID, update.MType)
	metrics := storage.Metrics[update.ID]
	if metrics == nil {
		storage.Logger.Infoln("Creating a new metric")
		storage.Metrics[update.ID] = &update
		return &update
	}
	switch update.MType {
	case entities.TypeGauge:
		storage.Logger.Infof("Old metric value: %f\n", *metrics.Value)
		metrics.Value = update.Value
		storage.Logger.Infof("New metric value: %f\n", *metrics.Value)
		return metrics
	case entities.TypeCounter:
		storage.Logger.Infof("Old metric value: %d\n", *metrics.Delta)
		var delta = *metrics.Delta + *update.Delta
		metrics.Delta = &delta
		storage.Logger.Infof("New metric value: %d\n", *metrics.Delta)
		return metrics
	}
	return nil // TODO: shouldn't happen, need to handle this
}

func (storage *MemStorage) Get(mType string, name string) (*entities.Metrics, error) {
	storage.Logger.Infof("Getting metric of type %s named %s\n", mType, name)
	result, ok := storage.Metrics[name]
	if !ok {
		storage.Logger.Infoln("No such metric")
		return nil, entities.ErrUnknownMetricName
	}
	switch mType {
	case entities.TypeGauge:
		fallthrough
	case entities.TypeCounter:
		if result.MType == mType {
			return result, nil
		}
	default:
		return nil, entities.ErrInvalidMetricType
	}
	return nil, entities.ErrUnknownMetricName
}

func (storage *MemStorage) GetAll() map[string]*entities.Metrics {
	return storage.Metrics
}
