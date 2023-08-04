package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/matthiasBT/monitoring/internal/storage"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetric(t *testing.T) {
	type args struct {
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
				url:      "/update/counter/Counter1/25",
				params:   []string{"counter", "Counter1", "25"},
				wantCode: http.StatusOK,
				storage:  emptyStorage(),
			},
		},
		{
			name: "correct gauge update",
			args: args{
				url:      "/update/gauge/Gauge1/25.4",
				params:   []string{"gauge", "Gauge1", "25.4"},
				wantCode: http.StatusOK,
				storage:  emptyStorage(),
			},
		},
		{
			name: "missing metric name",
			args: args{
				url:      "/update/counter//4",
				params:   []string{"counter", "", "4"},
				wantCode: http.StatusNotFound,
				storage:  emptyStorage(),
			},
		},
		{
			name: "invalid metric type",
			args: args{
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
			r := httptest.NewRequest(http.MethodPost, tt.args.url, nil)
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

func TestGetMetric(t *testing.T) {
	type args struct {
		url      string
		params   []string
		wantCode int
		wantBody []byte
		storage  *storage.MemStorage
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "get existing gauge",
			args: args{
				url:      "/value/gauge/Gauge1",
				params:   []string{"gauge", "Gauge1"},
				wantCode: http.StatusOK,
				wantBody: []byte("1.23"),
				storage:  nonEmptyStorage(),
			},
		},
		{
			name: "get existing counter",
			args: args{
				url:      "/value/counter/Counter1",
				params:   []string{"counter", "Counter1"},
				wantCode: http.StatusOK,
				wantBody: []byte("1"),
				storage:  nonEmptyStorage(),
			},
		},
		{
			name: "get non-existent gauge",
			args: args{
				url:      "/value/gauge/Gauge3",
				params:   []string{"gauge", "Gauge3"},
				wantCode: http.StatusNotFound,
				wantBody: nil,
				storage:  nonEmptyStorage(),
			},
		},
		{
			name: "get non-existent counter",
			args: args{
				url:      "/value/counter/Counter3",
				params:   []string{"counter", "Counter3"},
				wantCode: http.StatusNotFound,
				wantBody: nil,
				storage:  nonEmptyStorage(),
			},
		},
		{
			name: "get metric with invalid type",
			args: args{
				url:      "/value/hist/Gauge1",
				params:   []string{"hist", "Gauge1"},
				wantCode: http.StatusNotFound,
				wantBody: nil,
				storage:  nonEmptyStorage(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			r := httptest.NewRequest(http.MethodGet, tt.args.url, nil)
			w := httptest.NewRecorder()
			c := e.NewContext(r, w)
			c.SetPath("/value/:type/:name")
			c.SetParamNames("type", "name")
			c.SetParamValues(tt.args.params...)
			GetMetric(c, tt.args.storage)
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

func emptyStorage() *storage.MemStorage {
	return &storage.MemStorage{
		MetricsGauge:   make(map[string]float64),
		MetricsCounter: make(map[string]int64),
	}
}

func nonEmptyStorage() *storage.MemStorage {
	return &storage.MemStorage{
		MetricsGauge:   map[string]float64{"Gauge1": 1.23, "Gauge2": 1.49},
		MetricsCounter: map[string]int64{"Counter1": 1, "Counter2": 2},
	}
}

func TestGetAllMetrics(t *testing.T) {
	tests := []struct {
		name     string
		stor     *storage.MemStorage
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
		e := echo.New()
		e.Renderer = GetRenderer("../../web/template/*.html")
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		c := e.NewContext(r, w)
		c.SetPath("/")
		GetAllMetrics(c, tt.stor)
		res := w.Result()
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		if tt.wantBody != nil {
			assert.Equal(t, string(tt.wantBody), string(body))
		}
		assert.Equal(t, tt.wantCode, res.StatusCode)
	}
}
