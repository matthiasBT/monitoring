package adapters

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	HMACKey    []byte
}

func (r *HTTPReportAdapter) Report(metrics *common.Metrics) error {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	payload, err := json.Marshal(metrics)
	if err != nil {
		r.Logger.Errorf("Failed to marshal a metric: %v", metrics)
		return err
	}
	return r.report(&payload)
}

func (r *HTTPReportAdapter) ReportBatch(batch []*common.Metrics) error {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	payload, err := json.Marshal(batch)
	if err != nil {
		r.Logger.Errorf("Failed to marshal a batch of metrics: %v\n", err.Error())
		return err
	}
	return r.report(&payload)
}

func (r *HTTPReportAdapter) report(payload *[]byte) error {
	var (
		req *http.Request
		err error
	)
	u := url.URL{Scheme: "http", Host: r.ServerAddr, Path: r.UpdateURL}
	if req, err = r.createRequest(u, payload); err != nil {
		return err
	}
	if err := r.addHMACHeader(req, payload); err != nil {
		return err
	}

	f := func() (any, error) {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			r.Logger.Errorf("Request failed: %v\n", err.Error())
			return nil, err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	}
	bodyAny, err := r.Retrier.RetryChecked(context.Background(), f, utils.CheckConnectionError)
	if err != nil {
		return err
	}
	body := bodyAny.([]byte)
	r.Logger.Infof("Success. Server response: %v", string(body))
	return nil
}

func (r *HTTPReportAdapter) createRequest(path url.URL, payload *[]byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", path.String(), bytes.NewReader(*payload))
	if err != nil {
		r.Logger.Errorf("Failed to create a request: %v\n", err.Error())
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (r *HTTPReportAdapter) addHMACHeader(req *http.Request, payload *[]byte) error {
	if hash, err := r.hashData(payload); err != nil {
		return err
	} else if hash != "" {
		req.Header.Add("HashSHA256", hash)
	}
	return nil
}

func (r *HTTPReportAdapter) hashData(payload *[]byte) (string, error) {
	if bytes.Equal(r.HMACKey, []byte{}) {
		return "", nil
	}
	mac := hmac.New(sha256.New, r.HMACKey)
	if _, err := mac.Write(*payload); err != nil {
		r.Logger.Errorf("Failed to calculate hash: %v", err.Error())
		return "", err
	}
	hash := mac.Sum(nil)
	result := hex.EncodeToString(hash)
	r.Logger.Infof("HMAC-SHA256 hash: %s\n", result)
	return result, nil
}
