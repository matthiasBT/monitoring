// Package adapters provides functionality for managing in-memory storage
// of metrics. It includes operations for adding, retrieving, and flushing
// metrics, as well as periodic flush operations to an external storage.
package adapters

import (
	"context"
	"sync"
	"time"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/entities"
)

// State represents the in-memory storage state containing metrics
// and a mutex for synchronization.
type State struct {
	Metrics map[string]*common.Metrics // In-memory map of metrics
	Lock    *sync.Mutex                // Mutex for synchronization
}

// MemStorage is a struct that manages in-memory storage operations,
// periodic flushing to external storage, and logging.
type MemStorage struct {
	State                   // Embedded in-memory state
	Done   <-chan struct{}  // Channel signaling the end of the application
	Tick   <-chan time.Time // Ticker channel for periodic operations
	Logger logging.ILogger  // Logger for logging activities
	Keeper entities.Keeper  // External storage Keeper for flushing data
}

// NewMemStorage creates and returns a new MemStorage instance.
// It initializes in-memory state and sets up channels for periodic flushing.
func NewMemStorage(
	done <-chan struct{},
	tick <-chan time.Time,
	logger logging.ILogger,
	keeper entities.Keeper,
) entities.Storage {
	return &MemStorage{
		State: State{
			Metrics: make(map[string]*common.Metrics),
			Lock:    &sync.Mutex{},
		},
		Done:   done,
		Tick:   tick,
		Logger: logger,
		Keeper: keeper,
	}
}

// Add adds or updates a single metric in the in-memory storage.
func (storage *MemStorage) Add(ctx context.Context, update *common.Metrics) (*common.Metrics, error) {
	storage.Lock.Lock()
	defer storage.Lock.Unlock()

	return storage.addSingle(ctx, update)
}

// AddBatch adds a batch of metrics to the in-memory storage.
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

// Get retrieves a single metric from the in-memory storage based on query criteria.
func (storage *MemStorage) Get(ctx context.Context, query *common.Metrics) (*common.Metrics, error) {
	storage.Logger.Infof("Getting the metric %s %s\n", query.ID, query.MType)
	result, ok := storage.Metrics[query.ID]
	if !ok || result.MType != query.MType {
		storage.Logger.Errorf("No such metric\n")
		return nil, common.ErrUnknownMetric
	}
	return result, nil
}

// GetAll returns all metrics currently stored in-memory.
func (storage *MemStorage) GetAll(ctx context.Context) (map[string]*common.Metrics, error) {
	return storage.Metrics, nil
}

// Snapshot creates and returns a snapshot of the current in-memory metrics.
func (storage *MemStorage) Snapshot(context.Context) ([]*common.Metrics, error) {
	result := make([]*common.Metrics, 0, len(storage.Metrics))
	for _, val := range storage.Metrics {
		result = append(result, val)
	}
	return result, nil
}

// Init initializes the in-memory storage with provided data.
func (storage *MemStorage) Init(data []*common.Metrics) {
	storage.Logger.Infoln("Initializing the storage with new data. Old data will be lost")
	result := make(map[string]*common.Metrics, len(data))
	for _, metrics := range data {
		result[metrics.ID] = metrics
	}
	storage.Metrics = result
	storage.Logger.Infoln("Init finished successfully")
}

// Ping delegates the ping operation to the Keeper, if available.
func (storage *MemStorage) Ping(ctx context.Context) error {
	if storage.Keeper != nil {
		return storage.Keeper.Ping(ctx)
	}
	return nil
}

// FlushPeriodic handles periodic flushing of in-memory data to external storage.
func (storage *MemStorage) FlushPeriodic(ctx context.Context) {
	storage.Logger.Infoln("Launching the FlushPeriodic job")
	for {
		select {
		case <-storage.Done:
			storage.Logger.Infoln("Stopping the FlushPeriodic job")
			if err := storage.flush(ctx); err != nil {
				panic(err)
			}
			return
		case tick := <-storage.Tick:
			storage.Logger.Infof("The FlushPeriodic job is ticking at %v\n", tick)
			if err := storage.flush(ctx); err != nil {
				storage.Logger.Errorf("Failed to flush data: %s\n", err.Error())
			}
		}
	}
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
