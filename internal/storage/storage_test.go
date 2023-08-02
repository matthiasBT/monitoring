package storage

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemStorage_Add(t *testing.T) {
	tests := []struct {
		name         string
		update       []MetricUpdate
		wantGauges   map[string]float64
		wantCounters map[string]int64
	}{
		{
			name:         "add single counter",
			update:       []MetricUpdate{{TypeCounter, "Counter1", "1"}},
			wantGauges:   map[string]float64{},
			wantCounters: map[string]int64{"Counter1": 1},
		},
		{
			name: "add multiple counters",
			update: []MetricUpdate{
				{TypeCounter, "Counter1", "1"},
				{TypeCounter, "Counter2", "47"},
				{TypeCounter, "Counter1", "12"},
				{TypeCounter, "Counter2", "31"},
			},
			wantGauges:   map[string]float64{},
			wantCounters: map[string]int64{"Counter1": 13, "Counter2": 78},
		},
		{
			name:         "add single gauge",
			update:       []MetricUpdate{{TypeGauge, "Gauge1", "1.0"}},
			wantGauges:   map[string]float64{"Gauge1": 1.0},
			wantCounters: map[string]int64{},
		},
		{
			name: "add multiple gauges",
			update: []MetricUpdate{
				{TypeGauge, "Gauge1", "1.0"},
				{TypeGauge, "Gauge2", "5.43"},
				{TypeGauge, "Gauge1", "-33.11"},
				{TypeGauge, "Gauge3", "0.67"},
			},
			wantGauges:   map[string]float64{"Gauge1": -33.11, "Gauge2": 5.43, "Gauge3": 0.67},
			wantCounters: map[string]int64{},
		},
		{
			name: "mix gauges and counters",
			update: []MetricUpdate{
				{TypeGauge, "Gauge1", "1.0"},
				{TypeCounter, "Counter1", "5"},
				{TypeGauge, "Gauge2", "-33.11"},
				{TypeCounter, "Counter2", "22"},
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
			assert.InDeltaMapValues(
				t, tt.wantCounters, storage.MetricsCounter,
				0.0,
				"Counters don't match. Expected %v, got %v",
				tt.wantCounters,
				storage.MetricsCounter,
			)
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
			wantErr: ErrInvalidMetricVal,
		},
		{
			name:    "invalid metric type",
			fields:  fields{"hist", "Hist1", "0.567"},
			wantErr: ErrInvalidMetricType,
		},
		{
			name:    "missing metric name",
			fields:  fields{"counter", "", "4"},
			wantErr: ErrMissingMetricName,
		},
		{
			name:    "missing metric name with only whitespace chars",
			fields:  fields{"counter", "   \t\r\n\f   ", "4"},
			wantErr: ErrMissingMetricName,
		},
		{
			name:    "invalid counter value",
			fields:  fields{"counter", "Counter1", "four"},
			wantErr: ErrInvalidMetricVal,
		},
		{
			name:    "invalid gauge value",
			fields:  fields{"gauge", "Gauge1", "four-point-six"},
			wantErr: ErrInvalidMetricVal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MetricUpdate{
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
