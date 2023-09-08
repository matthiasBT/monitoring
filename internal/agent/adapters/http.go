package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
)

type HTTPReportAdapter struct {
	Logger     logging.ILogger
	ServerAddr string
	UpdateURL  string
	Retrier    utils.Retrier
	Lock       *sync.Mutex
}

func (r *HTTPReportAdapter) Report(metrics *common.Metrics) error {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	body, err := json.Marshal(metrics)
	if err != nil {
		r.Logger.Errorf("Failed to marshal a metric: %v", metrics)
		return err
	}

	u := url.URL{Scheme: "http", Host: r.ServerAddr, Path: r.UpdateURL}
	f := func() (any, error) {
		return http.Post(u.String(), "application/json", bytes.NewReader(body))
	}
	respAny, err := r.Retrier.RetryChecked(context.Background(), f, utils.CheckConnectionError)
	if err != nil {
		r.Logger.Errorf("Request failed: %v\n", err.Error())
		return err
	}
	resp := respAny.(*http.Response)
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r.Logger.Infof("Success: %v", string(body))
	return nil
}

func (r *HTTPReportAdapter) ReportBatch(batch []*common.Metrics) error {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	payload, err := json.Marshal(batch)
	if err != nil {
		r.Logger.Errorf("Failed to marshal a batch of metrics: %v\n", err.Error())
		return err
	}

	u := url.URL{Scheme: "http", Host: r.ServerAddr, Path: r.UpdateURL}
	f := func() (any, error) {
		return http.Post(u.String(), "application/json", bytes.NewReader(payload))
	}
	respAny, err := r.Retrier.RetryChecked(context.Background(), f, utils.CheckConnectionError)
	if err != nil {
		r.Logger.Errorf("Request failed: %v\n", err.Error())
		return err
	}
	resp := respAny.(*http.Response)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r.Logger.Infof("Success: %v", string(body))
	return nil
}
