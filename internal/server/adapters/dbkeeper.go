// Package adapters provides functionality for managing database connections,
// particularly for PostgreSQL, and includes methods for performing operations
// with retry logic and error handling.
package adapters

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/migrations"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
	"github.com/matthiasBT/monitoring/internal/server/entities"
)

// DBKeeper is a struct that manages database operations and holds a SQL database connection,
// a logger for logging operations, a retrier for handling retry logic, and a mutex for
// synchronizing operations.
type DBKeeper struct {
	DB      *sql.DB         // Database connection
	Lock    *sync.Mutex     // Mutex for synchronization
	Logger  logging.ILogger // Logger for logging activities
	Retrier utils.Retrier   // Retrier for retry logic
}

// OpenDB opens a database connection
func OpenDB(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	return db
}

// NewDBKeeper creates and returns a new DBKeeper instance with the provided database,
// logger, and retrier. It tests the database connection and runs migrations.
func NewDBKeeper(db *sql.DB, logger logging.ILogger, retrier utils.Retrier) entities.Keeper {
	keeper := DBKeeper{DB: db, Logger: logger, Retrier: retrier, Lock: &sync.Mutex{}}
	if err := keeper.Ping(context.Background()); err != nil {
		panic(err)
	}
	migrations.Migrate(db)
	return &keeper
}

// Flush saves a snapshot of storage data (metrics) to the database using transactions,
// with retry logic for transient errors.
func (dbk *DBKeeper) Flush(ctx context.Context, storageSnapshot []*common.Metrics) error {
	dbk.Logger.Infoln("Starting saving the storage data")

	dbk.Lock.Lock()
	defer dbk.Lock.Unlock()

	txOpt := sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	}

	f := func() (any, error) {
		return dbk.DB.BeginTx(ctx, &txOpt)
	}
	txAny, err := dbk.Retrier.RetryChecked(ctx, f, utils.CheckConnectionError)
	if err != nil {
		dbk.Logger.Errorf("Failed to open a transaction: %s\n", err.Error())
	}
	var tx = txAny.(*sql.Tx)
	//nolint:errcheck
	defer tx.Commit()

	for _, metrics := range storageSnapshot {
		if _, err := dbk.addSingle(ctx, tx, metrics); err != nil {
			dbk.Logger.Errorf("Failed to update a metric from snapshot: %s\n", err.Error())
			tx.Rollback()
			return err
		}
	}
	dbk.Logger.Infoln("Saving complete")
	return nil
}

// Restore fetches and returns all metrics from the database, with retry logic for transient errors.
func (dbk *DBKeeper) Restore() []*common.Metrics {
	dbk.Logger.Infoln("Restoring the storage data")
	var result []*common.Metrics
	ctx := context.Background()

	f := func() (any, error) {
		rows, err := dbk.DB.QueryContext(ctx, "SELECT * FROM metrics")
		if err != nil {
			return nil, err
		}
		// because of linter this check must be done twice - here and after the retrier
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return rows, nil
	}
	rowsAny, err := dbk.Retrier.RetryChecked(ctx, f, utils.CheckConnectionError)
	if err != nil {
		dbk.Logger.Errorf("Failed to fetch all table: %s\n", err.Error())
		panic(err)
	}

	var rows = rowsAny.(*sql.Rows)
	defer rows.Close()
	for rows.Next() {
		var metrics common.Metrics
		if err = scanSingleMetric(rows, &metrics); err != nil {
			dbk.Logger.Errorf("Failed to scan metric: %s\n", err.Error())
			panic(err)
		}
		result = append(result, &metrics)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	dbk.Logger.Infoln("Success")
	return result
}

func (dbk *DBKeeper) addSingle(ctx context.Context, tx *sql.Tx, update *common.Metrics) (*common.Metrics, error) {
	dbk.Logger.Infof("Updating a metric %s %s\n", update.ID, update.MType)

	metrics, err := dbk.get(ctx, tx, update)
	if err != nil {
		return nil, err
	}

	if metrics == nil {
		dbk.Logger.Infoln("Creating a new metric")
		if err := dbk.create(ctx, tx, update); err != nil {
			return nil, err
		}
		return update, nil
	}

	if result, err := dbk.update(ctx, tx, update); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (dbk *DBKeeper) get(ctx context.Context, tx *sql.Tx, search *common.Metrics) (*common.Metrics, error) {
	query := "SELECT * FROM metrics WHERE id = $1 AND mtype = $2"
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, query, search.ID, search.MType)
	} else {
		stmt, err := dbk.prepareStatement(ctx, query)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()
		row = stmt.QueryRowContext(ctx, search.ID, search.MType)
	}
	var result common.Metrics
	if err := scanMetric(row, &result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			dbk.Logger.Infof("No row found with ID %s and type %s\n", search.ID, search.MType)
			return nil, nil
		} else {
			dbk.Logger.Errorf("Failed to find metric %s %s\n", search.ID, search.MType, err.Error())
			return nil, err
		}
	}
	return &result, nil
}

func (dbk *DBKeeper) create(ctx context.Context, tx *sql.Tx, create *common.Metrics) error {
	query := `
		INSERT INTO metrics(id, mtype, delta, val)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET mtype = excluded.mtype, delta = excluded.delta, val = excluded.val
		WHERE metrics.id = excluded.id
		RETURNING *
	`
	var err error
	if tx == nil {
		f := func() (any, error) {
			return dbk.DB.ExecContext(ctx, query, create.ID, create.MType, create.Delta, create.Value)
		}
		_, err = dbk.Retrier.RetryChecked(ctx, f, utils.CheckConnectionError)
	} else {
		_, err = tx.ExecContext(ctx, query, create.ID, create.MType, create.Delta, create.Value)
	}
	if err != nil {
		dbk.Logger.Errorf("Failed to create a new metric %s\n", err.Error())
		return err
	}
	dbk.Logger.Infof("Created: %s %s\n", create.ID, create.MType)
	return nil
}

func (dbk *DBKeeper) update(ctx context.Context, tx *sql.Tx, update *common.Metrics) (*common.Metrics, error) {
	var row *sql.Row
	if update.MType == common.TypeCounter {
		query := "UPDATE metrics SET delta = delta + $1 WHERE id = $2 RETURNING *"
		if tx == nil {
			stmt, err := dbk.prepareStatement(ctx, query)
			if err != nil {
				return nil, err
			}
			defer stmt.Close()
			row = stmt.QueryRowContext(ctx, update.Delta, update.ID)
		} else {
			row = tx.QueryRowContext(ctx, query, update.Delta, update.ID)
		}
	} else {
		query := "UPDATE metrics SET val = $1 WHERE id = $2 RETURNING *"
		if tx == nil {
			stmt, err := dbk.prepareStatement(ctx, query)
			if err != nil {
				return nil, err
			}
			defer stmt.Close()
			row = stmt.QueryRowContext(ctx, update.Value, update.ID)
		} else {
			row = tx.QueryRowContext(ctx, query, update.Value, update.ID)
		}
	}

	var result common.Metrics
	if err := scanMetric(row, &result); err != nil {
		dbk.Logger.Errorf("Failed to update metric %s %s\n", update.ID, err.Error())
		return nil, err
	}
	dbk.Logger.Infof("Updated: %s %s\n", update.ID, update.MType)
	return &result, nil
}

func (dbk *DBKeeper) prepareStatement(ctx context.Context, query string) (*sql.Stmt, error) {
	f := func() (any, error) {
		return dbk.DB.PrepareContext(ctx, query)
	}
	stmtAny, err := dbk.Retrier.RetryChecked(ctx, f, utils.CheckConnectionError)
	if err != nil {
		dbk.Logger.Errorf("Failed to open a statement: %s\n", err.Error())
		return nil, err
	}
	return stmtAny.(*sql.Stmt), nil
}

func scanMetric(row *sql.Row, result *common.Metrics) error {
	return row.Scan(&result.ID, &result.MType, &result.Delta, &result.Value)
}

func scanSingleMetric(rows *sql.Rows, result *common.Metrics) error {
	return rows.Scan(&result.ID, &result.MType, &result.Delta, &result.Value)
}

// Shutdown closes the database connection and logs any errors encountered during the operation.
func (dbk *DBKeeper) Shutdown() {
	dbk.Logger.Infoln("Shutting down the database")
	if err := dbk.DB.Close(); err != nil {
		dbk.Logger.Errorf("Failed to shutdown the database: %s\n", err.Error())
		panic(err)
	}
}

// Ping tests the database connection and logs any errors.
func (dbk *DBKeeper) Ping(ctx context.Context) error {
	if err := dbk.DB.PingContext(ctx); err != nil {
		dbk.Logger.Errorf("Database ping failed: %s\n", err.Error())
		return err
	}
	return nil
}
