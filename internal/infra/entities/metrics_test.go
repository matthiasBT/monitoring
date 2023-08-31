package entities

import (
	"errors"
	"testing"
)

func TestMetrics_ValueAsString(t *testing.T) {
	tests := []struct {
		name string
		m    Metrics
		want string
	}{
		{
			name: "gauge",
			m: Metrics{
				ID:    "foobar",
				MType: TypeGauge,
				Delta: nil,
				Value: ptrfloat64(45.4235),
			},
			want: "45.4235",
		},
		{
			name: "counter",
			m: Metrics{
				ID:    "foobar",
				MType: TypeCounter,
				Delta: ptrint64(3015),
				Value: nil,
			},
			want: "3015",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.ValueAsString(); got != tt.want {
				t.Errorf("ValueAsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetrics_Validate(t *testing.T) {
	tests := []struct {
		name      string
		m         Metrics
		withValue bool
		wantErr   error
	}{
		{
			name: "get query",
			m: Metrics{
				ID:    "foobar",
				MType: TypeCounter,
				Delta: nil,
				Value: nil,
			},
			withValue: false,
			wantErr:   nil,
		},
		{
			name: "get query without name",
			m: Metrics{
				ID:    "",
				MType: TypeCounter,
				Delta: nil,
				Value: nil,
			},
			withValue: false,
			wantErr:   ErrMissingMetricName,
		},
		{
			name: "get query with incorrect type",
			m: Metrics{
				ID:    "foobar",
				MType: "hist",
				Delta: nil,
				Value: nil,
			},
			withValue: false,
			wantErr:   ErrInvalidMetricType,
		},
		{
			name: "update",
			m: Metrics{
				ID:    "foobar",
				MType: TypeCounter,
				Delta: ptrint64(123),
				Value: nil,
			},
			withValue: true,
			wantErr:   nil,
		},
		{
			name: "both values present",
			m: Metrics{
				ID:    "foobar",
				MType: TypeCounter,
				Delta: ptrint64(123),
				Value: ptrfloat64(3.45),
			},
			withValue: true,
			wantErr:   ErrInvalidMetricVal,
		},
		{
			name: "none values present",
			m: Metrics{
				ID:    "foobar",
				MType: TypeCounter,
				Delta: nil,
				Value: nil,
			},
			withValue: true,
			wantErr:   ErrInvalidMetricVal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.Validate(tt.withValue)
			if (err != nil && tt.wantErr == nil) ||
				(err == nil && tt.wantErr != nil) ||
				(err != nil && tt.wantErr != nil && !errors.Is(err, tt.wantErr)) {
				t.Errorf("Validate() error. got: %v, want: %v\n", err, tt.wantErr)
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
