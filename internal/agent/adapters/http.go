package adapters

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
)

type HTTPReportAdapter struct {
	Logger     logging.ILogger
	ServerAddr string
	UpdateURL  string
}

func (r *HTTPReportAdapter) Report(metrics *common.Metrics) error {
	body, err := json.Marshal(metrics)
	if err != nil {
		r.Logger.Infof("Failed to marshal a metric: %v", metrics)
		return err
	}

	u := url.URL{Scheme: "http", Host: r.ServerAddr, Path: r.UpdateURL}
	resp, err := http.Post(u.String(), "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r.Logger.Infof("Success: %v", string(body))
	return nil
}
