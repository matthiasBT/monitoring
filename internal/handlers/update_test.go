package handlers

import (
	"github.com/matthiasBT/monitoring/internal/storage"
)

const patternUpdate = "/update/"

//func TestUpdateMetric(t *testing.T) {
//	type args struct {
//		method   string
//		url      string
//		wantCode int
//		storage  *storage.MemStorage
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "get request",
//			args: args{
//				method:   http.MethodGet,
//				url:      "/update/counter/Counter1/25",
//				wantCode: http.StatusMethodNotAllowed,
//				storage:  emptyStorage(),
//			},
//		},
//		{
//			name: "correct counter update",
//			args: args{
//				method:   http.MethodPost,
//				url:      "/update/counter/Counter1/25",
//				wantCode: http.StatusOK,
//				storage:  emptyStorage(),
//			},
//		},
//		{
//			name: "correct gauge update",
//			args: args{
//				method:   http.MethodPost,
//				url:      "/update/gauge/Gauge/25.4",
//				wantCode: http.StatusOK,
//				storage:  emptyStorage(),
//			},
//		},
//		{
//			name: "empty update",
//			args: args{
//				method:   http.MethodPost,
//				url:      "/update/",
//				wantCode: http.StatusBadRequest,
//				storage:  emptyStorage(),
//			},
//		},
//		{
//			name: "missing metric name",
//			args: args{
//				method:   http.MethodPost,
//				url:      "/update/counter",
//				wantCode: http.StatusNotFound,
//				storage:  emptyStorage(),
//			},
//		},
//		{
//			name: "missing metric value",
//			args: args{
//				method:   http.MethodPost,
//				url:      "/update/counter/Counter1",
//				wantCode: http.StatusBadRequest,
//				storage:  emptyStorage(),
//			},
//		},
//		{
//			name: "invalid metric type",
//			args: args{
//				method:   http.MethodPost,
//				url:      "/update/hist/Counter1/4",
//				wantCode: http.StatusBadRequest,
//				storage:  emptyStorage(),
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := httptest.NewRequest(tt.args.method, tt.args.url, nil)
//			w := httptest.NewRecorder()
//			UpdateMetric(w, r, patternUpdate, tt.args.storage)
//			res := w.Result()
//			defer res.Body.Close()
//			assert.Equal(t, res.StatusCode, tt.args.wantCode)
//		})
//	}
//}

func emptyStorage() *storage.MemStorage {
	return &storage.MemStorage{
		MetricsGauge:   make(map[string]float64),
		MetricsCounter: make(map[string]int64),
	}
}
