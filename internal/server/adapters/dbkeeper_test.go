package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/matthiasBT/monitoring/internal/infra/config/server"
	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
	"github.com/matthiasBT/monitoring/internal/server/entities"
)

// todo
func TestDBKeeper_Flush(t *testing.T) {
	type fields struct {
		DB      *sql.DB
		Logger  logging.ILogger
		Retrier utils.Retrier
		Lock    *sync.Mutex
	}
	type args struct {
		ctx             context.Context
		storageSnapshot []*common.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbk := &DBKeeper{
				DB:      tt.fields.DB,
				Logger:  tt.fields.Logger,
				Retrier: tt.fields.Retrier,
				Lock:    tt.fields.Lock,
			}
			if err := dbk.Flush(tt.args.ctx, tt.args.storageSnapshot); (err != nil) != tt.wantErr {
				t.Errorf("Flush() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBKeeper_Ping(t *testing.T) {
	type fields struct {
		Retrier utils.Retrier
		Lock    *sync.Mutex
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "ping_success",
			fields: fields{
				Retrier: utils.Retrier{},
				Lock:    &sync.Mutex{},
			},
			args:    args{ctx: context.Background()},
			wantErr: nil,
		},
		{
			name: "ping_error",
			fields: fields{
				Retrier: utils.Retrier{},
				Lock:    &sync.Mutex{},
			},
			args:    args{ctx: context.Background()},
			wantErr: fmt.Errorf("fake error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("Error creating mock DB: %v", err)
			}
			defer db.Close()

			if tt.wantErr == nil {
				mock.ExpectPing()
			} else {
				mock.ExpectPing().WillReturnError(tt.wantErr)
			}
			dbk := &DBKeeper{
				DB:      db,
				Logger:  logging.SetupLogger(),
				Retrier: tt.fields.Retrier,
				Lock:    tt.fields.Lock,
			}
			if err := dbk.Ping(tt.args.ctx); !errors.Is(err, tt.wantErr) {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDBKeeper_Restore(t *testing.T) {
	type fields struct {
		Retrier utils.Retrier
		Lock    *sync.Mutex
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*common.Metrics
		wantErr error
	}{
		{
			name: "restore_success",
			fields: fields{
				Retrier: utils.Retrier{
					Attempts:         2,
					IntervalFirst:    100 * time.Millisecond,
					IntervalIncrease: 100 * time.Millisecond,
					Logger:           logging.SetupLogger(),
				},
				Lock: &sync.Mutex{},
			},
			want:    getMetricsRows(),
			wantErr: nil,
		},
		{
			name: "restore_query_failure",
			fields: fields{
				Retrier: utils.Retrier{
					Attempts:         2,
					IntervalFirst:    100 * time.Millisecond,
					IntervalIncrease: 100 * time.Millisecond,
					Logger:           logging.SetupLogger(),
				},
				Lock: &sync.Mutex{},
			},
			want:    nil,
			wantErr: fmt.Errorf("fake error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Error creating mock database: %v", err)
			}
			defer db.Close()
			rows := sqlmock.NewRows([]string{"ID", "MType", "Delta", "Value"}).
				AddRow("foo", "counter", "4", nil).
				AddRow("bar", "gauge", nil, "3.2")
			if tt.wantErr == nil {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM metrics")).WillReturnRows(rows)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM metrics")).WillReturnError(tt.wantErr)
			}
			dbk := &DBKeeper{
				DB:      db,
				Logger:  logging.SetupLogger(),
				Retrier: tt.fields.Retrier,
				Lock:    tt.fields.Lock,
			}
			defer func() {
				if r := recover(); r != nil {
					if tt.wantErr != nil {
						if !errors.Is(r.(error), tt.wantErr) {
							t.Errorf("Panic with unexpected error: %v", r)
						}
						if err := mock.ExpectationsWereMet(); err != nil {
							t.Errorf("Unfulfilled expectations: %s", err)
						}
					}
				}
			}()
			if got := dbk.Restore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Restore() = %v, want %v", got, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDBKeeper_Shutdown(t *testing.T) {
	type fields struct {
		Retrier utils.Retrier
		Lock    *sync.Mutex
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "shutdown_success",
			fields: fields{
				Retrier: utils.Retrier{},
				Lock:    &sync.Mutex{},
			},
			args:    args{ctx: context.Background()},
			wantErr: nil,
		},
		{
			name: "shutdown_error",
			fields: fields{
				Retrier: utils.Retrier{},
				Lock:    &sync.Mutex{},
			},
			args:    args{ctx: context.Background()},
			wantErr: fmt.Errorf("fake error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Error creating mock DB: %v", err)
			}
			defer db.Close()

			if tt.wantErr == nil {
				mock.ExpectClose()
			} else {
				mock.ExpectClose().WillReturnError(tt.wantErr)
			}
			dbk := &DBKeeper{
				DB:      db,
				Logger:  logging.SetupLogger(),
				Retrier: tt.fields.Retrier,
				Lock:    tt.fields.Lock,
			}
			defer func() {
				if r := recover(); r != nil {
					if tt.wantErr != nil {
						if !errors.Is(r.(error), tt.wantErr) {
							t.Errorf("Panic with unexpected error: %v", r)
						}
						if err := mock.ExpectationsWereMet(); err != nil {
							t.Errorf("Unfulfilled expectations: %s", err)
						}
					}
				}
			}()
			dbk.Shutdown()
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

// todo
func TestDBKeeper_addSingle(t *testing.T) {
	type fields struct {
		DB      *sql.DB
		Logger  logging.ILogger
		Retrier utils.Retrier
		Lock    *sync.Mutex
	}
	type args struct {
		ctx    context.Context
		tx     *sql.Tx
		update *common.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *common.Metrics
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbk := &DBKeeper{
				DB:      tt.fields.DB,
				Logger:  tt.fields.Logger,
				Retrier: tt.fields.Retrier,
				Lock:    tt.fields.Lock,
			}
			got, err := dbk.addSingle(tt.args.ctx, tt.args.tx, tt.args.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("addSingle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addSingle() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// todo
func TestDBKeeper_create(t *testing.T) {
	type fields struct {
		DB      *sql.DB
		Logger  logging.ILogger
		Retrier utils.Retrier
		Lock    *sync.Mutex
	}
	type args struct {
		ctx    context.Context
		tx     *sql.Tx
		create *common.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbk := &DBKeeper{
				DB:      tt.fields.DB,
				Logger:  tt.fields.Logger,
				Retrier: tt.fields.Retrier,
				Lock:    tt.fields.Lock,
			}
			if err := dbk.create(tt.args.ctx, tt.args.tx, tt.args.create); (err != nil) != tt.wantErr {
				t.Errorf("create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// todo
func TestDBKeeper_get(t *testing.T) {
	type fields struct {
		DB      *sql.DB
		Logger  logging.ILogger
		Retrier utils.Retrier
		Lock    *sync.Mutex
	}
	type args struct {
		ctx    context.Context
		tx     *sql.Tx
		search *common.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *common.Metrics
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbk := &DBKeeper{
				DB:      tt.fields.DB,
				Logger:  tt.fields.Logger,
				Retrier: tt.fields.Retrier,
				Lock:    tt.fields.Lock,
			}
			got, err := dbk.get(tt.args.ctx, tt.args.tx, tt.args.search)
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// todo
func TestDBKeeper_update(t *testing.T) {
	type fields struct {
		DB      *sql.DB
		Logger  logging.ILogger
		Retrier utils.Retrier
		Lock    *sync.Mutex
	}
	type args struct {
		ctx    context.Context
		tx     *sql.Tx
		update *common.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *common.Metrics
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbk := &DBKeeper{
				DB:      tt.fields.DB,
				Logger:  tt.fields.Logger,
				Retrier: tt.fields.Retrier,
				Lock:    tt.fields.Lock,
			}
			got, err := dbk.update(tt.args.ctx, tt.args.tx, tt.args.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// todo
func TestNewDBKeeper(t *testing.T) {
	type args struct {
		conf    *server.Config
		logger  logging.ILogger
		retrier utils.Retrier
	}
	tests := []struct {
		name string
		args args
		want entities.Keeper
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDBKeeper(tt.args.conf, tt.args.logger, tt.args.retrier); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDBKeeper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getMetricsRows() []*common.Metrics {
	counter := common.Metrics{
		ID:    "foo",
		MType: "counter",
		Delta: ptrint64(4),
		Value: nil,
	}
	gauge := common.Metrics{
		ID:    "bar",
		MType: "gauge",
		Delta: nil,
		Value: ptrfloat64(3.2),
	}
	return []*common.Metrics{&counter, &gauge}
}
