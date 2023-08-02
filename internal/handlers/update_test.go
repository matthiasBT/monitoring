package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/matthiasBT/monitoring/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetric(t *testing.T) {
	type args struct {
		method   string
		url      string
		params   []string
		wantCode int
		storage  *storage.MemStorage
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "correct counter update",
			args: args{
				method:   http.MethodPost,
				url:      "/update/counter/Counter1/25",
				params:   []string{"counter", "Counter1", "25"},
				wantCode: http.StatusOK,
				storage:  emptyStorage(),
			},
		},
		{
			name: "correct gauge update",
			args: args{
				method:   http.MethodPost,
				url:      "/update/gauge/Gauge1/25.4",
				params:   []string{"gauge", "Gauge1", "25.4"},
				wantCode: http.StatusOK,
				storage:  emptyStorage(),
			},
		},
		{
			name: "missing metric name",
			args: args{
				method:   http.MethodPost,
				url:      "/update/counter//4",
				params:   []string{"counter", "", "4"},
				wantCode: http.StatusNotFound,
				storage:  emptyStorage(),
			},
		},
		{
			name: "invalid metric type",
			args: args{
				method:   http.MethodPost,
				url:      "/update/hist/Hist1/4.879",
				params:   []string{"hist", "Hist1", "4.879"},
				wantCode: http.StatusBadRequest,
				storage:  emptyStorage(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			r := httptest.NewRequest(tt.args.method, tt.args.url, nil)
			w := httptest.NewRecorder()
			c := e.NewContext(r, w)
			c.SetPath("/update/:type/:name/:value")
			c.SetParamNames("type", "name", "value")
			c.SetParamValues(tt.args.params...)
			UpdateMetric(c, tt.args.storage)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.args.wantCode, res.StatusCode)
		})
	}
}

func emptyStorage() *storage.MemStorage {
	return &storage.MemStorage{
		MetricsGauge:   make(map[string]float64),
		MetricsCounter: make(map[string]int64),
	}
}
