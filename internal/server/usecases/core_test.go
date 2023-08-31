package usecases

import (
	"reflect"
	"testing"

	"github.com/matthiasBT/monitoring/internal/infra/entities"
)

func Test_prepareTemplateData(t *testing.T) {
	tests := []struct {
		name    string
		Metrics map[string]*entities.Metrics
		want    map[string]string
		wantErr error
	}{
		{
			name:    "get empty data for template",
			Metrics: make(map[string]*entities.Metrics),
			want:    make(map[string]string),
			wantErr: nil,
		},
		{
			name: "get mixed data for template",
			Metrics: map[string]*entities.Metrics{
				"FooBar": {
					ID:    "FooBar",
					MType: entities.TypeCounter,
					Delta: ptrint64(33),
					Value: nil,
				},
				"BarFoo": {
					ID:    "BarFoo",
					MType: entities.TypeGauge,
					Delta: nil,
					Value: ptrfloat64(55.1534),
				},
			},
			want: map[string]string{
				"FooBar": "33", "BarFoo": "55.1534",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := prepareTemplateData(tt.Metrics)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepareTemplateData() got = %v, want %v", got, tt.want)
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
