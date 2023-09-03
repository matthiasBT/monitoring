package adapters

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
)

type DBManager struct {
	db     *sql.DB
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
	d.db = db
	return nil
}

func (d *DBManager) Shutdown() {
	d.Logger.Infoln("Shutting down the database")
	if err := d.db.Close(); err != nil {
		d.Logger.Errorf("Failed to shutdown the database: %s\n", err.Error())
		panic(err)
	}
}

func (d *DBManager) Ping(ctx context.Context) error {
	if err := d.db.PingContext(ctx); err != nil {
		d.Logger.Errorf("Database ping failed: %s\n", err.Error())
		return err
	}
	return nil
}
