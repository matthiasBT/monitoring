package adapters

import (
	"errors"
	"sync"
	"testing"

	"github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
)

func TestMemStorage_Add(t *testing.T) {
	stor := MemStorage{
		Metrics: nil,
		Logger:  logging.SetupLogger(),
		Lock:    &sync.Mutex{},
	}
	tests := []struct {
		name        string
		Metrics     map[string]*entities.Metrics
		update      entities.Metrics
		wantMetrics map[string]*entities.Metrics
		want        entities.Metrics
		wantErr     error
	}{
		{
			name:    "create a counter",
			Metrics: make(map[string]*entities.Metrics),
			update: entities.Metrics{
				ID:    "FooBar",
				MType: entities.TypeCounter,
				Delta: ptrint64(33),
				Value: nil,
			},
			want: entities.Metrics{
				ID:    "FooBar",
				MType: entities.TypeCounter,
				Delta: ptrint64(33),
				Value: nil,
			},
			wantMetrics: map[string]*entities.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: entities.TypeCounter,
				Delta: ptrint64(33),
				Value: nil,
			}},
			wantErr: nil,
		},
		{
			name: "update a counter",
			Metrics: map[string]*entities.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: entities.TypeCounter,
				Delta: ptrint64(101),
				Value: nil,
			}},
			update: entities.Metrics{
				ID:    "FooBar",
				MType: entities.TypeCounter,
				Delta: ptrint64(99),
				Value: nil,
			},
			want: entities.Metrics{
				ID:    "FooBar",
				MType: entities.TypeCounter,
				Delta: ptrint64(200),
				Value: nil,
			},
			wantMetrics: map[string]*entities.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: entities.TypeCounter,
				Delta: ptrint64(200),
				Value: nil,
			}},
			wantErr: nil,
		},
		{
			name:    "create a gauge",
			Metrics: make(map[string]*entities.Metrics),
			update: entities.Metrics{
				ID:    "FooBar",
				MType: entities.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(44.1),
			},
			want: entities.Metrics{
				ID:    "FooBar",
				MType: entities.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(44.1),
			},
			wantMetrics: map[string]*entities.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: entities.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(44.1),
			}},
			wantErr: nil,
		},
		{
			name: "update a gauge",
			Metrics: map[string]*entities.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: entities.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(44.1),
			}},
			update: entities.Metrics{
				ID:    "FooBar",
				MType: entities.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(44.7),
			},
			want: entities.Metrics{
				ID:    "FooBar",
				MType: entities.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(44.7),
			},
			wantMetrics: map[string]*entities.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: entities.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(44.7),
			}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		stor.Metrics = tt.Metrics
		t.Run(tt.name, func(t *testing.T) {
			got, err := stor.Add(tt.update)
			if (err != nil && tt.wantErr == nil) ||
				(err == nil && tt.wantErr != nil) ||
				(err != nil && tt.wantErr != nil && !errors.Is(err, tt.wantErr)) {
				t.Errorf("Error mismatch. got: %v, want: %v\n", err, tt.wantErr)
				return
			}
			if !compare(got, &tt.want) {
				t.Errorf("Add() got = %v, want %v", got, tt.want)
			}
			if !compareState(stor.Metrics, tt.wantMetrics) {
				t.Errorf("State mismatch. got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_Get(t *testing.T) {
	stor := MemStorage{
		Metrics: nil,
		Logger:  logging.SetupLogger(),
		Lock:    &sync.Mutex{},
	}
	tests := []struct {
		name    string
		Metrics map[string]*entities.Metrics
		query   entities.Metrics
		want    *entities.Metrics
		wantErr error
	}{
		{
			name:    "get a counter from empty storage",
			Metrics: make(map[string]*entities.Metrics),
			query: entities.Metrics{
				ID:    "FooBar",
				MType: entities.TypeCounter,
			},
			want:    nil,
			wantErr: entities.ErrUnknownMetric,
		},
		{
			name: "get an existing counter",
			Metrics: map[string]*entities.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: entities.TypeCounter,
				Delta: ptrint64(33),
				Value: nil,
			}},
			query: entities.Metrics{
				ID:    "FooBar",
				MType: entities.TypeCounter,
			},
			want: &entities.Metrics{
				ID:    "FooBar",
				MType: entities.TypeCounter,
				Delta: ptrint64(33),
				Value: nil,
			},
			wantErr: nil,
		},
		{
			name:    "get a gauge from empty storage",
			Metrics: make(map[string]*entities.Metrics),
			query: entities.Metrics{
				ID:    "FooBar",
				MType: entities.TypeGauge,
			},
			want:    nil,
			wantErr: entities.ErrUnknownMetric,
		},
		{
			name: "get an existing gauge",
			Metrics: map[string]*entities.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: entities.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(77.1),
			}},
			query: entities.Metrics{
				ID:    "FooBar",
				MType: entities.TypeGauge,
			},
			want: &entities.Metrics{
				ID:    "FooBar",
				MType: entities.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(77.1),
			},
			wantErr: nil,
		},
		{
			name: "name clash",
			Metrics: map[string]*entities.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: entities.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(77.1),
			}},
			query: entities.Metrics{
				ID:    "FooBar",
				MType: entities.TypeCounter,
			},
			want:    nil,
			wantErr: entities.ErrUnknownMetric,
		},
	}
	for _, tt := range tests {
		stor.Metrics = tt.Metrics
		t.Run(tt.name, func(t *testing.T) {
			got, err := stor.Get(tt.query)
			if (err != nil && tt.wantErr == nil) ||
				(err == nil && tt.wantErr != nil) ||
				(err != nil && tt.wantErr != nil && !errors.Is(err, tt.wantErr)) {
				t.Errorf("Error mismatch. got: %v, want: %v\n", err, tt.wantErr)
				return
			}
			if !compare(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func ptrfloat64(val float64) *float64 {
	return &val
}

func ptrint64(val int64) *int64 {
	return &val
}

func compare(m1 *entities.Metrics, m2 *entities.Metrics) bool {
	return m1 == nil && m2 == nil ||
		m1.ID == m2.ID &&
			m1.MType == m2.MType &&
			(m1.Delta != nil && m2.Delta != nil && *m1.Delta == *m2.Delta ||
				m1.Delta == nil && m2.Delta == nil) &&
			(m1.Value != nil && m2.Value != nil && *m1.Value == *m2.Value ||
				m1.Value == nil && m2.Value == nil)
}

func compareState(got map[string]*entities.Metrics, want map[string]*entities.Metrics) bool {
	for key, m1 := range got {
		if !compare(m1, want[key]) {
			return false
		}
	}
	for key, m2 := range want {
		if !compare(m2, got[key]) {
			return false
		}
	}
	return true
}
