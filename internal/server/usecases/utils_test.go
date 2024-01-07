package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
	common "github.com/matthiasBT/monitoring/internal/infra/entities"
)

func Test_handleInvalidMetric(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
		wantBody string
	}{
		{
			name:     "invalid_type",
			err:      common.ErrInvalidMetricType,
			wantCode: http.StatusBadRequest,
			wantBody: common.ErrInvalidMetricType.Error(),
		},
		{
			name:     "no_name",
			err:      common.ErrMissingMetricName,
			wantCode: http.StatusNotFound,
			wantBody: common.ErrMissingMetricName.Error(),
		},
		{
			name:     "invalid_val",
			err:      common.ErrInvalidMetricVal,
			wantCode: http.StatusBadRequest,
			wantBody: common.ErrInvalidMetricVal.Error(),
		},
		{
			name:     "internal_server_error",
			err:      errors.New("foobar"),
			wantCode: http.StatusInternalServerError,
			wantBody: "foobar",
		},
	}
	for _, tt := range tests {
		rr := httptest.NewRecorder()
		t.Run(tt.name, func(t *testing.T) {
			handleInvalidMetric(rr, tt.err)
		})
		if tt.wantCode != rr.Code {
			t.Errorf("Code mismatch. got: %v, want: %v\n", rr.Code, tt.wantCode)
		}
		if !bytes.Equal([]byte(tt.wantBody), rr.Body.Bytes()) {
			t.Errorf("Body mismatch. got: %s, want: %s\n", rr.Body.Bytes(), tt.wantBody)
		}
	}
}

func Test_parseMetric(t *testing.T) {
	type args struct {
		url        string
		body       string
		asJSON     bool
		withValue  bool
		paramType  string
		paramName  string
		paramValue string
	}
	tests := []struct {
		name string
		args args
		want *common.Metrics
	}{
		{
			name: "valid_counter_json", args: args{
				url:       "/update/",
				body:      `{"id": "foobar", "type": "counter", "delta": 123}`,
				asJSON:    true,
				withValue: false,
			}, want: &common.Metrics{
				ID:    "foobar",
				MType: common.TypeCounter,
				Delta: ptrint64(123),
				Value: nil,
			},
		},
		{
			name: "valid_gauge_json", args: args{
				url:       "/update/",
				body:      `{"id": "foobar", "type": "gauge", "value": 123.4}`,
				asJSON:    true,
				withValue: false,
			}, want: &common.Metrics{
				ID:    "foobar",
				MType: common.TypeGauge,
				Delta: nil,
				Value: ptrfloat64(123.4),
			},
		},
		{
			name: "invalid_counter_json", args: args{
				url:       "/update/",
				body:      `"id": "foobar", "type": "counter", "delta": 123}`,
				asJSON:    true,
				withValue: false,
			}, want: nil,
		},
		{
			name: "valid_counter_params_noval", args: args{
				url:        "/update/counter/foobar/123",
				body:       "",
				asJSON:     false,
				withValue:  false,
				paramType:  common.TypeCounter,
				paramName:  "foobar",
				paramValue: "123",
			}, want: &common.Metrics{
				ID:    "foobar",
				MType: common.TypeCounter,
				Delta: nil,
				Value: nil,
			},
		},
		{
			name: "valid_gauge_params_noval", args: args{
				url:        "/update/gauge/foobar/123",
				body:       "",
				asJSON:     false,
				withValue:  false,
				paramType:  "gauge",
				paramName:  "foobar",
				paramValue: "123.4",
			}, want: &common.Metrics{
				ID:    "foobar",
				MType: "gauge",
				Delta: nil,
				Value: nil,
			},
		},
		{
			name: "valid_counter_params_val", args: args{
				url:        "/update/counter/foobar/123",
				body:       "",
				asJSON:     false,
				withValue:  true,
				paramType:  common.TypeCounter,
				paramName:  "foobar",
				paramValue: "123",
			}, want: &common.Metrics{
				ID:    "foobar",
				MType: common.TypeCounter,
				Delta: ptrint64(123),
				Value: nil,
			},
		},
		{
			name: "valid_gauge_params_val", args: args{
				url:        "/update/gauge/foobar/123.4",
				body:       "",
				asJSON:     false,
				withValue:  true,
				paramType:  common.TypeGauge,
				paramName:  "foobar",
				paramValue: "123.4",
			}, want: &common.Metrics{
				ID:    "foobar",
				MType: "gauge",
				Delta: nil,
				Value: ptrfloat64(123.4),
			},
		},
		{
			name: "invalid_counter_params_val", args: args{
				url:        "/update/counter/foobar/123a",
				body:       "",
				asJSON:     false,
				withValue:  true,
				paramType:  common.TypeCounter,
				paramName:  "foobar",
				paramValue: "123a",
			}, want: nil,
		},
		{
			name: "invalid_gauge_params_val", args: args{
				url:        "/update/gauge/foobar/123.4a",
				body:       "",
				asJSON:     false,
				withValue:  true,
				paramType:  common.TypeGauge,
				paramName:  "foobar",
				paramValue: "123.4a",
			}, want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error
			if tt.args.asJSON {
				body := []byte(tt.args.body)
				req, err = http.NewRequest("_", tt.args.url, bytes.NewBuffer(body))
				if err != nil {
					t.Fatal(err)
				}
			} else {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("type", tt.args.paramType)
				rctx.URLParams.Add("name", tt.args.paramName)
				rctx.URLParams.Add("value", tt.args.paramValue)
				req, err = http.NewRequest("_", tt.args.url, bytes.NewBuffer([]byte{}))
				if err != nil {
					t.Fatal(err)
				}
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			}
			if got := parseMetric(req, tt.args.asJSON, tt.args.withValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeMetric(t *testing.T) {
	type args struct {
		asJSON  bool
		metrics common.Metrics
	}
	tests := []struct {
		name     string
		args     args
		wantBody string
	}{
		{
			name: "valid_counter_json",
			args: args{
				asJSON: true,
				metrics: common.Metrics{
					ID:    "foobar",
					MType: common.TypeCounter,
					Delta: ptrint64(123),
					Value: nil,
				},
			},
			wantBody: `{"id":"foobar","type":"counter","delta":123}`,
		},
		{
			name: "valid_gauge_json",
			args: args{
				asJSON: true,
				metrics: common.Metrics{
					ID:    "foobar",
					MType: common.TypeGauge,
					Delta: nil,
					Value: ptrfloat64(123.4),
				},
			},
			wantBody: `{"id":"foobar","type":"gauge","value":123.4}`,
		},
		{
			name: "valid_counter_plain",
			args: args{
				asJSON: false,
				metrics: common.Metrics{
					ID:    "foobar",
					MType: common.TypeCounter,
					Delta: ptrint64(123),
					Value: nil,
				},
			},
			wantBody: "123",
		},
		{
			name: "valid_gauge_plain",
			args: args{
				asJSON: false,
				metrics: common.Metrics{
					ID:    "foobar",
					MType: common.TypeGauge,
					Delta: nil,
					Value: ptrfloat64(123.4),
				},
			},
			wantBody: "123.4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			writeMetric(rr, tt.args.asJSON, &tt.args.metrics)
			var want, got common.Metrics
			json.Unmarshal([]byte(tt.wantBody), &want)
			json.Unmarshal(rr.Body.Bytes(), &got)
			if !reflect.DeepEqual(want, got) {
				t.Errorf("Body mismatch. got: %v, want: %v\n", got, want)
			}
		})
	}
}
