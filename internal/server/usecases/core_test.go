package usecases

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/adapters"
	"github.com/matthiasBT/monitoring/internal/server/entities"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetric(t *testing.T) {
	logger := logging.SetupLogger()
	type args struct {
		params   map[string]string
		wantCode int
		storage  entities.Storage
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "correct counter update",
			args: args{
				params:   map[string]string{"type": "counter", "name": "Counter1", "value": "25"},
				wantCode: http.StatusOK,
				storage:  emptyStorage(),
			},
		},
		{
			name: "correct gauge update",
			args: args{
				params:   map[string]string{"type": "gauge", "name": "Gauge1", "value": "25.4"},
				wantCode: http.StatusOK,
				storage:  emptyStorage(),
			},
		},
		{
			name: "missing metric name",
			args: args{
				params:   map[string]string{"type": "counter", "name": "", "value": "4"},
				wantCode: http.StatusNotFound,
				storage:  emptyStorage(),
			},
		},
		{
			name: "invalid metric type",
			args: args{
				params:   map[string]string{"type": "hist", "name": "Hist1", "value": "4.879"},
				wantCode: http.StatusBadRequest,
				storage:  emptyStorage(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := &BaseController{Stor: tt.args.storage, Logger: logger}
			w := httptest.NewRecorder()
			UpdateMetric(w, controller, tt.args.params)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.args.wantCode, res.StatusCode)
		})
	}
}

func TestGetMetric(t *testing.T) {
	logger := logging.SetupLogger()
	type args struct {
		params   map[string]string
		wantCode int
		wantBody []byte
		storage  entities.Storage
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "get existing gauge",
			args: args{
				params:   map[string]string{"type": "gauge", "name": "Gauge1"},
				wantCode: http.StatusOK,
				wantBody: []byte("1.23"),
				storage:  nonEmptyStorage(),
			},
		},
		{
			name: "get existing counter",
			args: args{
				params:   map[string]string{"type": "counter", "name": "Counter1"},
				wantCode: http.StatusOK,
				wantBody: []byte("1"),
				storage:  nonEmptyStorage(),
			},
		},
		{
			name: "get non-existent gauge",
			args: args{
				params:   map[string]string{"type": "gauge", "name": "Gauge3"},
				wantCode: http.StatusNotFound,
				wantBody: nil,
				storage:  nonEmptyStorage(),
			},
		},
		{
			name: "get non-existent counter",
			args: args{
				params:   map[string]string{"type": "counter", "name": "Counter3"},
				wantCode: http.StatusNotFound,
				wantBody: nil,
				storage:  nonEmptyStorage(),
			},
		},
		{
			name: "get metric with invalid type",
			args: args{
				params:   map[string]string{"type": "hist", "name": "Hist1"},
				wantCode: http.StatusNotFound,
				wantBody: nil,
				storage:  nonEmptyStorage(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			controller := &BaseController{Stor: tt.args.storage, Logger: logger}
			GetMetric(w, controller, tt.args.params)
			res := w.Result()
			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)
			if tt.args.wantBody != nil {
				assert.Equal(t, body, tt.args.wantBody)
			}
			assert.Equal(t, tt.args.wantCode, res.StatusCode)
		})
	}
}

func emptyStorage() entities.Storage {
	return &adapters.MemStorage{
		MetricsGauge:   make(map[string]float64),
		MetricsCounter: make(map[string]int64),
		Logger:         logging.SetupLogger(),
	}
}

func nonEmptyStorage() entities.Storage {
	return &adapters.MemStorage{
		MetricsGauge:   map[string]float64{"Gauge1": 1.23, "Gauge2": 1.49},
		MetricsCounter: map[string]int64{"Counter1": 1, "Counter2": 2},
		Logger:         logging.SetupLogger(),
	}
}

func TestGetAllMetrics(t *testing.T) {
	logger := logging.SetupLogger()
	tests := []struct {
		name     string
		stor     entities.Storage
		wantBody []byte
		wantErr  error
		wantCode int
	}{
		{
			name: "empty metrics table",
			stor: emptyStorage(),
			wantBody: []byte(`<html>
<body>
<h1>Metrics table</h1>
<ul>
    
</ul>
</body>
</html>`),
			wantCode: http.StatusOK,
			wantErr:  nil,
		},
		{
			name: "non-empty metrics table",
			stor: nonEmptyStorage(),
			wantBody: []byte(`<html>
<body>
<h1>Metrics table</h1>
<ul>
    <li> Counter1 : 1 </li>
    <li> Counter2 : 2 </li>
    <li> Gauge1 : 1.23 </li>
    <li> Gauge2 : 1.49 </li>
    
</ul>
</body>
</html>`),
			wantCode: http.StatusOK,
			wantErr:  nil,
		},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		controller := &BaseController{Stor: tt.stor, TemplatePath: "../../../web/template", Logger: logger}
		GetAllMetrics(w, controller, "all_metrics.html")
		res := w.Result()
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		if tt.wantBody != nil {
			assert.Equal(t, string(tt.wantBody), string(body))
		}
		assert.Equal(t, tt.wantCode, res.StatusCode)
	}
}
