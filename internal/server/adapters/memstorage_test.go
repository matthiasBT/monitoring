package adapters

import (
	"context"
	"errors"
	"sync"
	"testing"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
)

type FakeKeeper struct {
	calledFlush    bool
	calledRestore  bool
	calledPing     bool
	calledShutdown bool
}

func (k *FakeKeeper) Flush(context.Context, []*common.Metrics) error {
	k.calledFlush = true
	return nil
}

func (k *FakeKeeper) Restore() []*common.Metrics {
	k.calledRestore = true
	return nil
}

func (k *FakeKeeper) Ping(ctx context.Context) error {
	k.calledPing = true
	return nil
}

func (k *FakeKeeper) Shutdown() {
	k.calledShutdown = true
}

func TestMemStorage_Add(t *testing.T) {
	st := State{
		Metrics: nil,
		Lock:    &sync.Mutex{},
	}
	stor := MemStorage{
		State:  st,
		Logger: logging.SetupLogger(),
	}
	tests := []struct {
		name        string
		Metrics     map[string]*common.Metrics
		update      common.Metrics
		wantMetrics map[string]*common.Metrics
		want        common.Metrics
		wantErr     error
	}{
		{
			name:    "create a counter",
			Metrics: make(map[string]*common.Metrics),
			update: common.Metrics{
				ID:    "FooBar",
				MType: common.TypeCounter,
				Delta: ptrint64(33),
				Value: nil,
			},
			want: common.Metrics{
				ID:    "FooBar",
				MType: common.TypeCounter,
				Delta: ptrint64(33),
				Value: nil,
			},
			wantMetrics: map[string]*common.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: common.TypeCounter,
				Delta: ptrint64(33),
				Value: nil,
			}},
			wantErr: nil,
		},
		{
			name: "update a counter",
			Metrics: map[string]*common.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: common.TypeCounter,
				Delta: ptrint64(101),
				Value: nil,
			}},
			update: common.Metrics{
				ID:    "FooBar",
				MType: common.TypeCounter,
				Delta: ptrint64(99),
				Value: nil,
			},
			want: common.Metrics{
				ID:    "FooBar",
				MType: common.TypeCounter,
				Delta: ptrint64(200),
				Value: nil,
			},
			wantMetrics: map[string]*common.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: common.TypeCounter,
				Delta: ptrint64(200),
				Value: nil,
			}},
			wantErr: nil,
		},
		{
			name:    "create a gauge",
			Metrics: make(map[string]*common.Metrics),
			update: common.Metrics{
				ID:    "FooBar",
				MType: common.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(44.1),
			},
			want: common.Metrics{
				ID:    "FooBar",
				MType: common.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(44.1),
			},
			wantMetrics: map[string]*common.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: common.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(44.1),
			}},
			wantErr: nil,
		},
		{
			name: "update a gauge",
			Metrics: map[string]*common.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: common.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(44.1),
			}},
			update: common.Metrics{
				ID:    "FooBar",
				MType: common.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(44.7),
			},
			want: common.Metrics{
				ID:    "FooBar",
				MType: common.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(44.7),
			},
			wantMetrics: map[string]*common.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: common.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(44.7),
			}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		stor.Metrics = tt.Metrics
		t.Run(tt.name, func(t *testing.T) {
			got, err := stor.Add(context.Background(), &tt.update)
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

func TestMemStorage_AddBatch(t *testing.T) {
	st := State{
		Metrics: nil,
		Lock:    &sync.Mutex{},
	}
	stor := MemStorage{
		State:  st,
		Logger: logging.SetupLogger(),
	}
	tests := []struct {
		name        string
		Metrics     map[string]*common.Metrics
		update      []*common.Metrics
		wantMetrics map[string]*common.Metrics
	}{
		{
			name:    "create_counters",
			Metrics: make(map[string]*common.Metrics),
			update: []*common.Metrics{
				{
					ID:    "FooBar1",
					MType: common.TypeCounter,
					Delta: ptrint64(33),
					Value: nil,
				},
				{
					ID:    "FooBar2",
					MType: common.TypeCounter,
					Delta: ptrint64(66),
					Value: nil,
				},
			},
			wantMetrics: map[string]*common.Metrics{
				"FooBar1": {
					ID:    "FooBar1",
					MType: common.TypeCounter,
					Delta: ptrint64(33),
					Value: nil,
				},
				"FooBar2": {
					ID:    "FooBar2",
					MType: common.TypeCounter,
					Delta: ptrint64(66),
					Value: nil,
				},
			},
		},
		{
			name:    "create_gauges",
			Metrics: make(map[string]*common.Metrics),
			update: []*common.Metrics{
				{
					ID:    "BarFoo1",
					MType: common.TypeGauge,
					Delta: nil,
					Value: ptrfloat64(44.1),
				},
				{
					ID:    "BarFoo2",
					MType: common.TypeGauge,
					Delta: nil,
					Value: ptrfloat64(55.5),
				},
			},
			wantMetrics: map[string]*common.Metrics{
				"BarFoo1": {
					ID:    "BarFoo1",
					MType: common.TypeGauge,
					Delta: nil,
					Value: ptrfloat64(44.1),
				},
				"BarFoo2": {
					ID:    "BarFoo2",
					MType: common.TypeGauge,
					Delta: nil,
					Value: ptrfloat64(55.5),
				},
			},
		},
	}
	for _, tt := range tests {
		stor.Metrics = tt.Metrics
		t.Run(tt.name, func(t *testing.T) {
			stor.AddBatch(context.Background(), tt.update)
		})
		if !compareState(stor.Metrics, tt.wantMetrics) {
			t.Errorf("State mismatch. got = %v, want %v", stor.Metrics, tt.wantMetrics)
		}
	}
}

func TestMemStorage_Get(t *testing.T) {
	st := State{
		Metrics: nil,
		Lock:    &sync.Mutex{},
	}
	stor := MemStorage{
		State:  st,
		Logger: logging.SetupLogger(),
	}
	tests := []struct {
		name    string
		Metrics map[string]*common.Metrics
		query   common.Metrics
		want    *common.Metrics
		wantErr error
	}{
		{
			name:    "get a counter from empty storage",
			Metrics: make(map[string]*common.Metrics),
			query: common.Metrics{
				ID:    "FooBar",
				MType: common.TypeCounter,
			},
			want:    nil,
			wantErr: common.ErrUnknownMetric,
		},
		{
			name: "get an existing counter",
			Metrics: map[string]*common.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: common.TypeCounter,
				Delta: ptrint64(33),
				Value: nil,
			}},
			query: common.Metrics{
				ID:    "FooBar",
				MType: common.TypeCounter,
			},
			want: &common.Metrics{
				ID:    "FooBar",
				MType: common.TypeCounter,
				Delta: ptrint64(33),
				Value: nil,
			},
			wantErr: nil,
		},
		{
			name:    "get a gauge from empty storage",
			Metrics: make(map[string]*common.Metrics),
			query: common.Metrics{
				ID:    "FooBar",
				MType: common.TypeGauge,
			},
			want:    nil,
			wantErr: common.ErrUnknownMetric,
		},
		{
			name: "get an existing gauge",
			Metrics: map[string]*common.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: common.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(77.1),
			}},
			query: common.Metrics{
				ID:    "FooBar",
				MType: common.TypeGauge,
			},
			want: &common.Metrics{
				ID:    "FooBar",
				MType: common.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(77.1),
			},
			wantErr: nil,
		},
		{
			name: "name clash",
			Metrics: map[string]*common.Metrics{"FooBar": {
				ID:    "FooBar",
				MType: common.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(77.1),
			}},
			query: common.Metrics{
				ID:    "FooBar",
				MType: common.TypeCounter,
			},
			want:    nil,
			wantErr: common.ErrUnknownMetric,
		},
	}
	for _, tt := range tests {
		stor.Metrics = tt.Metrics
		t.Run(tt.name, func(t *testing.T) {
			got, err := stor.Get(context.Background(), &tt.query)
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

func compare(m1 *common.Metrics, m2 *common.Metrics) bool {
	return m1 == nil && m2 == nil ||
		m1.ID == m2.ID &&
			m1.MType == m2.MType &&
			(m1.Delta != nil && m2.Delta != nil && *m1.Delta == *m2.Delta ||
				m1.Delta == nil && m2.Delta == nil) &&
			(m1.Value != nil && m2.Value != nil && *m1.Value == *m2.Value ||
				m1.Value == nil && m2.Value == nil)
}

func compareState(got map[string]*common.Metrics, want map[string]*common.Metrics) bool {
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

func TestMemStorage_Init(t *testing.T) {
	tests := []struct {
		name string
		data []*common.Metrics
		want map[string]*common.Metrics
	}{
		{
			name: "new_state_after_init",
			data: []*common.Metrics{
				{
					ID:    "FooBar1",
					MType: common.TypeCounter,
					Delta: ptrint64(33),
					Value: nil,
				},
				{
					ID:    "FooBar2",
					MType: common.TypeCounter,
					Delta: ptrint64(66),
					Value: nil,
				},
			},
			want: map[string]*common.Metrics{
				"FooBar1": {
					ID:    "FooBar1",
					MType: common.TypeCounter,
					Delta: ptrint64(33),
					Value: nil,
				},
				"FooBar2": {
					ID:    "FooBar2",
					MType: common.TypeCounter,
					Delta: ptrint64(66),
					Value: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{Logger: logging.SetupLogger()}
			storage.Init(tt.data)
			if !compareState(storage.Metrics, tt.want) {
				t.Errorf("State mismatch. Got: %v, want: %v", storage.Metrics, tt.want)
			}
		})
	}
}

func TestMemStorage_Ping(t *testing.T) {
	tests := []struct {
		name   string
		Keeper *FakeKeeper
	}{
		{
			name:   "no_keeper_is_ok",
			Keeper: nil,
		},
		{
			name:   "keeper_ping_is_called",
			Keeper: &FakeKeeper{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{}
			if tt.Keeper != nil {
				storage.Keeper = tt.Keeper
			}
			storage.Ping(context.Background())
			if tt.Keeper != nil {
				if !tt.Keeper.calledPing ||
					tt.Keeper.calledFlush ||
					tt.Keeper.calledRestore ||
					tt.Keeper.calledShutdown {
					t.Errorf("Invalid interaction with the keeper. Keeper state: %v", tt.Keeper)
				}
			}
		})
	}
}

func TestMemStorage_Snapshot(t *testing.T) {
	tests := []struct {
		name  string
		state map[string]*common.Metrics
		want1 []*common.Metrics
		want2 []*common.Metrics
	}{
		{
			name: "snapshot_values_are_in_any_order",
			state: map[string]*common.Metrics{
				"FooBar1": {
					ID:    "FooBar1",
					MType: common.TypeCounter,
					Delta: ptrint64(33),
					Value: nil,
				},
				"FooBar2": {
					ID:    "FooBar2",
					MType: common.TypeCounter,
					Delta: ptrint64(66),
					Value: nil,
				},
			},
			want1: []*common.Metrics{
				{
					ID:    "FooBar1",
					MType: common.TypeCounter,
					Delta: ptrint64(33),
					Value: nil,
				},
				{
					ID:    "FooBar2",
					MType: common.TypeCounter,
					Delta: ptrint64(66),
					Value: nil,
				},
			},
			want2: []*common.Metrics{
				{
					ID:    "FooBar2",
					MType: common.TypeCounter,
					Delta: ptrint64(66),
					Value: nil,
				},
				{
					ID:    "FooBar1",
					MType: common.TypeCounter,
					Delta: ptrint64(33),
					Value: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{}
			storage.Metrics = tt.state
			got, _ := storage.Snapshot(context.Background())
			if len(got) != len(tt.state) {
				t.Errorf("Snapshot must contain %d elements, got: %d", len(tt.state), len(got))
			}
			if !(compare(got[0], tt.want1[0]) && compare(got[1], tt.want1[1])) &&
				!(compare(got[0], tt.want2[0]) && compare(got[1], tt.want2[1])) {
				t.Errorf("Snapshot() got = %v, want %v or %v", got, tt.want1, tt.want2)
			}
		})
	}
}
