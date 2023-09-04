package adapters

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/entities"
)

type DBStorage struct {
	DB     *sql.DB
	Lock   *sync.Mutex
	Logger logging.ILogger
	Keeper entities.Keeper
}

func NewDBStorage(db *sql.DB, logger logging.ILogger, keeper entities.Keeper) entities.Storage {
	return &DBStorage{
		Lock:   &sync.Mutex{},
		DB:     db,
		Logger: logger,
		Keeper: keeper,
	}
}

func (storage *DBStorage) SetKeeper(entities.Keeper) {
	storage.Logger.Errorf("No keeper is necessary for DBStorage")
}

func (storage *DBStorage) Add(ctx context.Context, update common.Metrics) (*common.Metrics, error) {
	storage.Lock.Lock()
	defer storage.Lock.Unlock()

	storage.Logger.Infof("Updating a metric %s %s\n", update.ID, update.MType)

	metrics, err := storage.get(ctx, &update)
	if err != nil {
		return nil, err
	}

	if metrics == nil || metrics.MType != update.MType {
		storage.Logger.Infoln("Creating a new metric")
		if err := storage.create(ctx, &update); err != nil {
			return nil, err
		}
		return &update, nil
	}

	if result, err := storage.update(ctx, &update); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (storage *DBStorage) Get(ctx context.Context, search common.Metrics) (*common.Metrics, error) {
	if metrics, err := storage.get(ctx, &search); err != nil {
		return nil, err
	} else if metrics == nil {
		return nil, common.ErrUnknownMetric
	} else {
		return metrics, nil
	}
}

func (storage *DBStorage) GetAll(ctx context.Context) (map[string]*common.Metrics, error) {
	query := "SELECT * FROM metrics"
	rows, err := storage.DB.QueryContext(ctx, query)
	if err != nil {
		storage.Logger.Errorf("Failed to fetch all table: %s\n", err.Error())
		return nil, err
	}

	defer rows.Close()
	var result = make(map[string]*common.Metrics)
	for rows.Next() {
		var metrics common.Metrics
		if err = scanSingleMetric(rows, &metrics); err != nil {
			storage.Logger.Errorf("Failed to scan metric: %s\n", err.Error())
			return nil, err
		}
		result[metrics.ID] = &metrics
	}
	return result, nil
}

func (storage *DBStorage) Snapshot(context.Context) ([]*common.Metrics, error) {
	storage.Logger.Errorf("No snapshot can be taken from DBStorage") // TODO: warn
	return nil, nil                                                  // TODO: return error?
}

func (storage *DBStorage) Init([]*common.Metrics) {
	storage.Logger.Errorf("No init is necessary for DBStorage") // TODO: warn
}

func (storage *DBStorage) flush(context.Context) {
	storage.Logger.Errorf("No flush is necessary for DBStorage") // TODO: warn
}

func (storage *DBStorage) get(ctx context.Context, search *common.Metrics) (*common.Metrics, error) {
	query := "SELECT * FROM metrics WHERE id = $1 AND mtype = $2"
	stmt, err := storage.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRowContext(ctx, search.ID, search.MType)

	var result common.Metrics
	if err := scanMetric(row, &result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			storage.Logger.Infof("No row found with ID and type %s\n", search.ID, search.MType)
			return nil, nil
		} else {
			storage.Logger.Errorf("Failed to find metric %s %s\n", search.ID, search.MType, err.Error())
			return nil, err
		}
	}
	return &result, nil
}

func (storage *DBStorage) create(ctx context.Context, create *common.Metrics) error {
	query := "INSERT INTO metrics(id, mtype, delta, val) VALUES ($1, $2, $3, $4)"
	_, err := storage.DB.ExecContext(ctx, query, create.ID, create.MType, create.Delta, create.Value)
	if err != nil {
		storage.Logger.Errorf("Failed to create a new metric %s\n", err.Error())
		return err
	}
	return nil
}

func (storage *DBStorage) update(ctx context.Context, update *common.Metrics) (*common.Metrics, error) {
	var row *sql.Row
	if update.MType == common.TypeCounter {
		query := "UPDATE metrics SET delta = delta + $1 WHERE id = $2 RETURNING *"
		stmt, err := storage.DB.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()
		row = stmt.QueryRowContext(ctx, update.Delta, update.ID)
	} else {
		query := "UPDATE metrics SET val = $1 WHERE id = $2 RETURNING *"
		stmt, err := storage.DB.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()
		row = stmt.QueryRowContext(ctx, update.Value, update.ID)
	}

	var result common.Metrics
	if err := scanMetric(row, &result); err != nil {
		storage.Logger.Errorf("Failed to update metric %s %s\n", update.ID, err.Error())
		return nil, err
	}
	return &result, nil
}

func scanMetric(row *sql.Row, result *common.Metrics) error {
	return row.Scan(&result.ID, &result.MType, &result.Delta, &result.Value)
}

func scanSingleMetric(rows *sql.Rows, result *common.Metrics) error {
	return rows.Scan(&result.ID, &result.MType, &result.Delta, &result.Value)
}
