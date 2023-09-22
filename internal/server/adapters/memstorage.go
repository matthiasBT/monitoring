package adapters

import (
	"context"
	"sync"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/entities"
	"golang.org/x/exp/maps"
)

type State struct {
	Metrics map[string]*common.Metrics
	Lock    *sync.Mutex
}

type MemStorage struct {
	State
	Logger logging.ILogger
	Keeper entities.Keeper
}

func NewMemStorage(logger logging.ILogger, keeper entities.Keeper) entities.Storage {
	return &MemStorage{
		State: State{
			Metrics: make(map[string]*common.Metrics),
			Lock:    &sync.Mutex{},
		},
		Logger: logger,
		Keeper: keeper,
	}
}

func (storage *MemStorage) SetKeeper(keeper entities.Keeper) {
	storage.Keeper = keeper
}

func (storage *MemStorage) Add(ctx context.Context, update *common.Metrics) (*common.Metrics, error) {
	storage.Lock.Lock()
	defer storage.Lock.Unlock()

	return storage.addSingle(ctx, update)
}

func (storage *MemStorage) AddBatch(ctx context.Context, batch []*common.Metrics) error {
	storage.Lock.Lock()
	defer storage.Lock.Unlock()

	for _, metrics := range batch {
		if _, err := storage.addSingle(ctx, metrics); err != nil {
			storage.Logger.Errorf("Failed to add metric from batch: %s\n", err.Error())
			return err
		}
	}
	return nil
}

func (storage *MemStorage) Get(ctx context.Context, query *common.Metrics) (*common.Metrics, error) {
	storage.Logger.Infof("Getting the metric %s %s\n", query.ID, query.MType)
	result, ok := storage.Metrics[query.ID]
	if !ok || result.MType != query.MType {
		storage.Logger.Errorf("No such metric\n")
		return nil, common.ErrUnknownMetric
	}
	return result, nil
}

func (storage *MemStorage) GetAll(ctx context.Context) (map[string]*common.Metrics, error) {
	return storage.Metrics, nil
}

func (storage *MemStorage) Snapshot(ctx context.Context) ([]*common.Metrics, error) {
	data := maps.Values(storage.Metrics)
	return data, nil
}

func (storage *MemStorage) Init(data []*common.Metrics) {
	storage.Logger.Infoln("Initializing the storage with new data. Old data will be lost")
	result := make(map[string]*common.Metrics, len(data))
	for _, metrics := range data {
		result[metrics.ID] = metrics
	}
	storage.Metrics = result
	storage.Logger.Infoln("Init finished successfully")
}

func (storage *MemStorage) Ping(ctx context.Context) error {
	if storage.Keeper != nil {
		return storage.Keeper.Ping(ctx)
	}
	return nil
}

func (storage *MemStorage) addSingle(ctx context.Context, update *common.Metrics) (*common.Metrics, error) {
	storage.Logger.Infof("Updating a metric %s %s\n", update.ID, update.MType)
	metrics := storage.Metrics[update.ID]
	if metrics == nil || metrics.MType != update.MType {
		storage.Logger.Infoln("Creating a new metric")
		storage.Metrics[update.ID] = update
		if err := storage.flush(ctx); err != nil {
			return nil, err
		}
		return update, nil
	}
	if update.MType == common.TypeGauge {
		storage.Logger.Infof("Old metric value: %f\n", *metrics.Value)
		metrics.Value = update.Value
		storage.Logger.Infof("New metric value: %f\n", *metrics.Value)
		if err := storage.flush(ctx); err != nil {
			return nil, err
		}
		return metrics, nil
	} else { // Counter
		storage.Logger.Infof("Old metric value: %d\n", *metrics.Delta)
		var delta = *metrics.Delta + *update.Delta
		metrics.Delta = &delta
		storage.Logger.Infof("New metric value: %d\n", *metrics.Delta)
		if err := storage.flush(ctx); err != nil {
			return nil, err
		}
		return metrics, nil
	}
}

func (storage *MemStorage) flush(ctx context.Context) error {
	if storage.Keeper != nil {
		snapshot, _ := storage.Snapshot(ctx)
		return storage.Keeper.Flush(ctx, snapshot)
	}
	return nil
}
