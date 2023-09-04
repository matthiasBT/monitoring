package adapters

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
)

const (
	createType = `
		DO $$ BEGIN
			CREATE TYPE metric_type AS ENUM('gauge', 'counter');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`
	createTable = `
		CREATE TABLE IF NOT EXISTS metrics (
			id text primary key,
			mtype metric_type,
			delta bigint,
			val double precision
		)
	`
	createIndex = `
		CREATE INDEX IF NOT EXISTS search_idx ON metrics
		USING btree(id, mtype)
	`
)

type DBManager struct {
	DB     *sql.DB
	Logger logging.ILogger
}

func NewDBManager(dsn string, logger logging.ILogger) (*DBManager, error) {
	d := DBManager{Logger: logger}
	return &d, d.Init(dsn)
}

func (d *DBManager) Init(dsn string) error {
	d.Logger.Infoln("Opening the database: %s\n", dsn) // TODO: switch to debug
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		d.Logger.Errorf("Failed to open the database: %s\n", err.Error())
		return err
	}
	d.DB = db
	d.Logger.Infoln("Trying to ping the database")
	return d.prepare()
}

func (d *DBManager) prepare() error {
	d.Logger.Infoln("Creating database objects if necessary")
	for _, query := range []string{createType, createTable, createIndex} {
		if _, err := d.DB.Exec(query); err != nil {
			d.Logger.Errorf("Failed to execute query: %s\n", err.Error())
			return err
		}
	}
	return nil
}

func (d *DBManager) Shutdown() {
	d.Logger.Infoln("Shutting down the database")
	if err := d.DB.Close(); err != nil {
		d.Logger.Errorf("Failed to shutdown the database: %s\n", err.Error())
		panic(err)
	}
}

func (d *DBManager) Ping(ctx context.Context) error {
	if err := d.DB.PingContext(ctx); err != nil {
		d.Logger.Errorf("Database ping failed: %s\n", err.Error())
		return err
	}
	return nil
}
