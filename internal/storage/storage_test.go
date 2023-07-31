package storage

import (
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
				t, storage.MetricsGauge,
				tt.wantGauges,
				0.0,
				"Gauges don't match. Got %v, want %v",
				storage.MetricsGauge,
				tt.wantGauges,
			)
			assert.InDeltaMapValues(
				t, storage.MetricsCounter,
				tt.wantCounters,
				0.0,
				"Counters don't match. Got %v, want %v",
				storage.MetricsCounter,
				tt.wantCounters,
			)
		})
	}
}
