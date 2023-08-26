package adapters

import (
	"errors"
	"testing"

	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/entities"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_Add(t *testing.T) {
	logger := logging.SetupLogger()
	tests := []struct {
		name         string
		update       []entities.MetricUpdate
		wantGauges   map[string]float64
		wantCounters map[string]int64
	}{
		{
			name:         "add single counter",
			update:       []entities.MetricUpdate{{entities.TypeCounter, "Counter1", "1"}},
			wantGauges:   map[string]float64{},
			wantCounters: map[string]int64{"Counter1": 1},
		},
		{
			name: "add multiple counters",
			update: []entities.MetricUpdate{
				{entities.TypeCounter, "Counter1", "1"},
				{entities.TypeCounter, "Counter2", "47"},
				{entities.TypeCounter, "Counter1", "12"},
				{entities.TypeCounter, "Counter2", "31"},
			},
			wantGauges:   map[string]float64{},
			wantCounters: map[string]int64{"Counter1": 13, "Counter2": 78},
		},
		{
			name:         "add single gauge",
			update:       []entities.MetricUpdate{{entities.TypeGauge, "Gauge1", "1.0"}},
			wantGauges:   map[string]float64{"Gauge1": 1.0},
			wantCounters: map[string]int64{},
		},
		{
			name: "add multiple gauges",
			update: []entities.MetricUpdate{
				{entities.TypeGauge, "Gauge1", "1.0"},
				{entities.TypeGauge, "Gauge2", "5.43"},
				{entities.TypeGauge, "Gauge1", "-33.11"},
				{entities.TypeGauge, "Gauge3", "0.67"},
			},
			wantGauges:   map[string]float64{"Gauge1": -33.11, "Gauge2": 5.43, "Gauge3": 0.67},
			wantCounters: map[string]int64{},
		},
		{
			name: "mix gauges and counters",
			update: []entities.MetricUpdate{
				{entities.TypeGauge, "Gauge1", "1.0"},
				{entities.TypeCounter, "Counter1", "5"},
				{entities.TypeGauge, "Gauge2", "-33.11"},
				{entities.TypeCounter, "Counter2", "22"},
			},
			wantGauges:   map[string]float64{"Gauge1": 1.0, "Gauge2": -33.11},
			wantCounters: map[string]int64{"Counter1": 5, "Counter2": 22},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{
				MetricsGauge:   make(map[string]float64),
				MetricsCounter: make(map[string]int64),
				Logger:         logger,
			}
			for _, upd := range tt.update {
				storage.Add(upd)
			}
			assert.InDeltaMapValues(
				t, tt.wantGauges, storage.MetricsGauge,
				0.0,
				"Gauges don't match. Expected %v, got %v",
				tt.wantGauges,
				storage.MetricsGauge,
			)
			assert.Equalf(t, tt.wantCounters, storage.MetricsCounter, "Counters don't match")
		})
	}
}

func TestMetricUpdate_Validate(t *testing.T) {
	type fields struct {
		Type  string
		Name  string
		Value string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name:    "correct gauge update",
			fields:  fields{"gauge", "Gauge1", "123.4"},
			wantErr: nil,
		},
		{
			name:    "correct counter update",
			fields:  fields{"counter", "Counter1", "123"},
			wantErr: nil,
		},
		{
			name:    "incorrect counter update with float",
			fields:  fields{"counter", "Counter1", "123.4"},
			wantErr: entities.ErrInvalidMetricVal,
		},
		{
			name:    "invalid metric type",
			fields:  fields{"hist", "Hist1", "0.567"},
			wantErr: entities.ErrInvalidMetricType,
		},
		{
			name:    "missing metric name",
			fields:  fields{"counter", "", "4"},
			wantErr: entities.ErrMissingMetricName,
		},
		{
			name:    "missing metric name with only whitespace chars",
			fields:  fields{"counter", "   \t\r\n\f   ", "4"},
			wantErr: entities.ErrMissingMetricName,
		},
		{
			name:    "invalid counter value",
			fields:  fields{"counter", "Counter1", "four"},
			wantErr: entities.ErrInvalidMetricVal,
		},
		{
			name:    "invalid gauge value",
			fields:  fields{"gauge", "Gauge1", "four-point-six"},
			wantErr: entities.ErrInvalidMetricVal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := entities.MetricUpdate{
				Type:  tt.fields.Type,
				Name:  tt.fields.Name,
				Value: tt.fields.Value,
			}
			err := m.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MetricUpdate.Validate error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMemStorage_Get(t *testing.T) {
	logger := logging.SetupLogger()
	type args struct {
		mType string
		name  string
	}
	tests := []struct {
		name    string
		fields  MemStorage
		args    args
		want    string
		wantErr error
	}{
		{
			name: "get existing gauge",
			fields: MemStorage{
				MetricsGauge:   map[string]float64{"Gauge1": 0.123},
				MetricsCounter: map[string]int64{},
				Logger:         logger,
			},
			args: args{
				mType: "gauge",
				name:  "Gauge1",
			},
			want:    "0.123",
			wantErr: nil,
		},
		{
			name: "get non-existent gauge",
			fields: MemStorage{
				MetricsGauge:   map[string]float64{"Gauge1": 0.123},
				MetricsCounter: map[string]int64{},
				Logger:         logger,
			},
			args: args{
				mType: "gauge",
				name:  "GaugeX",
			},
			want:    "",
			wantErr: entities.ErrUnknownMetricName,
		},
		{
			name: "get existing counter",
			fields: MemStorage{
				MetricsGauge:   map[string]float64{},
				MetricsCounter: map[string]int64{"Counter1": 56},
				Logger:         logger,
			},
			args: args{
				mType: "counter",
				name:  "Counter1",
			},
			want:    "56",
			wantErr: nil,
		},
		{
			name: "get non-existent counter",
			fields: MemStorage{
				MetricsGauge:   map[string]float64{},
				MetricsCounter: map[string]int64{"Counter1": 56},
				Logger:         logger,
			},
			args: args{
				mType: "counter",
				name:  "CounterX",
			},
			want:    "",
			wantErr: entities.ErrUnknownMetricName,
		},
		{
			name: "get same name metric",
			fields: MemStorage{
				MetricsGauge:   map[string]float64{"MemoryAlloc": 45.123},
				MetricsCounter: map[string]int64{"MemoryAlloc": 56},
				Logger:         logger,
			},
			args: args{
				mType: "counter",
				name:  "MemoryAlloc",
			},
			want:    "56",
			wantErr: nil,
		},
		{
			name: "get metric with incorrect type",
			fields: MemStorage{
				MetricsGauge:   map[string]float64{"SomeGauge": 45.123},
				MetricsCounter: map[string]int64{"SomeCounter": 56},
				Logger:         logger,
			},
			args: args{
				mType: "hist",
				name:  "SomeGauge",
			},
			want:    "",
			wantErr: entities.ErrInvalidMetricType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{
				MetricsGauge:   tt.fields.MetricsGauge,
				MetricsCounter: tt.fields.MetricsCounter,
				Logger:         tt.fields.Logger,
			}
			got, err := storage.Get(tt.args.mType, tt.args.name)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MetricUpdate.Validate error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equalf(t, tt.want, got, "Get(%v, %v)", tt.args.mType, tt.args.name)
		})
	}
}

func TestMemStorage_GetAll(t *testing.T) {
	logger := logging.SetupLogger()
	tests := []struct {
		name   string
		fields MemStorage
		want   map[string]string
	}{
		{
			name: "empty map",
			fields: MemStorage{
				MetricsGauge:   make(map[string]float64),
				MetricsCounter: make(map[string]int64),
				Logger:         logger,
			},
			want: make(map[string]string),
		},
		{
			name: "only counters map",
			fields: MemStorage{
				MetricsGauge:   make(map[string]float64),
				MetricsCounter: map[string]int64{"Counter1": 1, "Counter2": 2},
				Logger:         logger,
			},
			want: map[string]string{"Counter1": "1", "Counter2": "2"},
		},
		{
			name: "only gauges map",
			fields: MemStorage{
				MetricsGauge:   map[string]float64{"Gauge1": 1.234, "Gauge2": 2.345},
				MetricsCounter: make(map[string]int64),
				Logger:         logger,
			},
			want: map[string]string{"Gauge1": "1.234", "Gauge2": "2.345"},
		},
		{
			name: "both metric types map",
			fields: MemStorage{
				MetricsGauge:   map[string]float64{"Gauge1": 1.234, "Gauge2": 2.345},
				MetricsCounter: map[string]int64{"Counter1": 1, "Counter2": 2},
				Logger:         logger,
			},
			want: map[string]string{"Gauge1": "1.234", "Gauge2": "2.345", "Counter1": "1", "Counter2": "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{
				MetricsGauge:   tt.fields.MetricsGauge,
				MetricsCounter: tt.fields.MetricsCounter,
				Logger:         logger,
			}
			res := storage.GetAll()
			assert.Equalf(t, tt.want, res, "Storage contents don't match. Expected: %v, got: %v", tt.want, res)
		})
	}
}
