package handlers

import (
	"errors"
	"github.com/matthiasBT/monitoring/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

const prefix = "/update/"

func Test_parseMetricUpdate(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    *storage.MetricUpdate
		wantErr error
	}{
		{
			name:    "correct counter update",
			url:     "/update/counter/Counter1/31",
			want:    &storage.MetricUpdate{Type: storage.TypeCounter, Name: "Counter1", Value: "31"},
			wantErr: nil,
		},
		{
			name:    "correct gauge update",
			url:     "/update/gauge/Gauge1/31.2",
			want:    &storage.MetricUpdate{Type: storage.TypeGauge, Name: "Gauge1", Value: "31.2"},
			wantErr: nil,
		},
		{
			name:    "empty update",
			url:     "/update/",
			want:    nil,
			wantErr: ErrInvalidMetricType,
		},
		{
			name:    "missing metric name",
			url:     "/update/counter",
			want:    nil,
			wantErr: ErrMissingMetricName,
		},
		{
			name:    "missing metric value",
			url:     "/update/counter/Counter1",
			want:    nil,
			wantErr: ErrInvalidMetricVal,
		},
		{
			name:    "invalid metric value",
			url:     "/update/counter/Counter1/4/1",
			want:    nil,
			wantErr: ErrInvalidMetricVal,
		},
		{
			name:    "invalid metric type",
			url:     "/update/hist/Counter1/4",
			want:    nil,
			wantErr: ErrInvalidMetricType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMetricUpdate(tt.url, prefix)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("parseMetricUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.EqualValues(t, tt.want, got, "Updates don't match. Expected %v, got %v", tt.want, got)
		})
	}
}
