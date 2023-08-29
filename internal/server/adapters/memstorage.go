package adapters

import (
	"sync"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/entities"
)

type MemStorage struct {
	Metrics map[string]*common.Metrics
	Logger  logging.ILogger
	Keeper  entities.Keeper
	Lock    *sync.Mutex
}

func NewMemStorage(logger logging.ILogger, keeper entities.Keeper) entities.Storage {
	return &MemStorage{
		Metrics: make(map[string]*common.Metrics),
		Logger:  logger,
		Keeper:  keeper,
		Lock:    &sync.Mutex{},
	}
}

func (storage *MemStorage) SetKeeper(keeper entities.Keeper) {
	storage.Keeper = keeper
}

func (storage *MemStorage) Add(update common.Metrics) (*common.Metrics, error) {
	storage.Lock.Lock()
	defer storage.Lock.Unlock()

	storage.Logger.Infof("Updating a metric %s %s\n", update.ID, update.MType)
	metrics := storage.Metrics[update.ID]
	if metrics == nil || metrics.MType != update.MType {
		storage.Logger.Infoln("Creating a new metric")
		storage.Metrics[update.ID] = &update
		storage.flush()
		return &update, nil
	}
	if update.MType == common.TypeGauge {
		storage.Logger.Infof("Old metric value: %f\n", *metrics.Value)
		metrics.Value = update.Value
		storage.Logger.Infof("New metric value: %f\n", *metrics.Value)
		storage.flush()
		return metrics, nil
	} else { // Counter
		storage.Logger.Infof("Old metric value: %d\n", *metrics.Delta)
		var delta = *metrics.Delta + *update.Delta
		metrics.Delta = &delta
		storage.Logger.Infof("New metric value: %d\n", *metrics.Delta)
		storage.flush()
		return metrics, nil
	}
}

func (storage *MemStorage) Get(query common.Metrics) (*common.Metrics, error) {
	storage.Logger.Infof("Getting the metric %s %s\n", query.ID, query.MType)

	result, ok := storage.Metrics[query.ID]
	if !ok || result.MType != query.MType {
		storage.Logger.Errorf("No such metric\n")
		return nil, common.ErrUnknownMetric
	}
	return result, nil
}

func (storage *MemStorage) GetAll() (map[string]*common.Metrics, error) {
	return storage.Metrics, nil
}

func (storage *MemStorage) Init(state map[string]*common.Metrics) {
	storage.Metrics = state
}

func (storage *MemStorage) flush() {
	if storage.Keeper != nil {
		storage.Keeper.Flush()
	}
}
