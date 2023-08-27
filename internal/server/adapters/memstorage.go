package adapters

import (
	"sync"

	"github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
)

type MemStorage struct {
	Metrics map[string]*entities.Metrics
	Logger  logging.ILogger
	Events  chan<- struct{}
	Lock    *sync.Mutex
}

func (storage *MemStorage) Add(update entities.Metrics) (*entities.Metrics, error) {
	storage.Lock.Lock()
	defer storage.Lock.Unlock()

	storage.Logger.Infof("Updating a metric %s %s\n", update.ID, update.MType)
	metrics := storage.Metrics[update.ID]
	if metrics == nil || metrics.MType != update.MType {
		storage.Logger.Infoln("Creating a new metric")
		storage.Metrics[update.ID] = &update
		storage.sendEvent()
		return &update, nil
	}
	if update.MType == entities.TypeGauge {
		storage.Logger.Infof("Old metric value: %f\n", *metrics.Value)
		metrics.Value = update.Value
		storage.Logger.Infof("New metric value: %f\n", *metrics.Value)
		storage.sendEvent()
		return metrics, nil
	} else { // Counter
		storage.Logger.Infof("Old metric value: %d\n", *metrics.Delta)
		var delta = *metrics.Delta + *update.Delta
		metrics.Delta = &delta
		storage.Logger.Infof("New metric value: %d\n", *metrics.Delta)
		storage.sendEvent()
		return metrics, nil
	}
}

func (storage *MemStorage) Get(query entities.Metrics) (*entities.Metrics, error) {
	storage.Logger.Infof("Getting the metric %s %s\n", query.ID, query.MType)

	result, ok := storage.Metrics[query.ID]
	if !ok || result.MType != query.MType {
		storage.Logger.Errorf("No such metric\n")
		return nil, entities.ErrUnknownMetric
	}
	return result, nil
}

func (storage *MemStorage) GetAll() (map[string]*entities.Metrics, error) {
	return storage.Metrics, nil
}

func (storage *MemStorage) Init(state map[string]*entities.Metrics) {
	storage.Metrics = state
}

func (storage *MemStorage) sendEvent() {
	storage.Events <- struct{}{}
}
