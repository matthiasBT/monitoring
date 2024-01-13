// Package adapters provides functionalities to handle HTTP communication,
// specifically for reporting metrics. It includes structures and methods
// for sending reports, handling retries, and ensuring data integrity.
package adapters

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
)

// HTTPReportAdapter is responsible for sending metrics over HTTP. It handles
// the creation and sending of HTTP requests, including retries and HMAC authentication.
// It uses a channel to queue payloads for asynchronous processing.
type HTTPReportAdapter struct {
	// Logger is used for logging messages related to HTTP reporting activities.
	Logger logging.ILogger

	// Jobs is an internal channel used to queue payloads for reporting.
	Jobs chan []byte

	// ServerAddr specifies the HTTP server address where reports are sent.
	ServerAddr string

	// UpdateURL is the endpoint for updating or sending reports.
	UpdateURL string

	// HMACKey is the key used for HMAC-SHA256 hashing to ensure data integrity.
	HMACKey []byte

	// CryptoKey is the key used for payload encryption
	CryptoKey *rsa.PublicKey

	// Retrier is used to handle retries for HTTP requests in case of failures.
	Retrier utils.Retrier
}

// NewHTTPReportAdapter creates and returns a new HTTPReportAdapter. It initializes
// the adapter with the provided logger, server address, update URL, retrier, HMAC key,
// and sets up worker goroutines based on the provided workerNum.
func NewHTTPReportAdapter(
	logger logging.ILogger,
	serverAddr string,
	updateURL string,
	retrier utils.Retrier,
	hmacKey []byte,
	cryptoKey *rsa.PublicKey,
	workerNum uint,
) *HTTPReportAdapter {
	adapter := HTTPReportAdapter{
		Logger:     logger,
		ServerAddr: serverAddr,
		UpdateURL:  updateURL,
		Retrier:    retrier,
		HMACKey:    hmacKey,
		CryptoKey:  cryptoKey,
		Jobs:       make(chan []byte, workerNum),
	}
	var i uint
	for i = 0; i < workerNum; i++ {
		go func() {
			for {
				data := <-adapter.Jobs
				//nolint:errcheck
				adapter.report(data)
			}
		}()
	}
	return &adapter
}

// Report sends a single metric over HTTP. It marshals the metric, logs any errors
// in marshaling, and queues the payload for processing.
func (r *HTTPReportAdapter) Report(metrics *common.Metrics) error {
	payload, err := json.Marshal(metrics)
	if err != nil {
		r.Logger.Errorf("Failed to marshal a metric: %v", metrics)
		return err
	}
	r.Jobs <- payload
	return nil
}

// ReportBatch sends a batch of metrics over HTTP. It marshals the batch of metrics,
// logs any errors in marshaling, and queues the payload for processing.
func (r *HTTPReportAdapter) ReportBatch(batch []*common.Metrics) error {
	payload, err := json.Marshal(batch)
	if err != nil {
		r.Logger.Errorf("Failed to marshal a batch of metrics: %v\n", err.Error())
		return err
	}
	r.Jobs <- payload
	return nil
}

func (r *HTTPReportAdapter) report(payload []byte) error {
	var (
		req *http.Request
		err error
	)
	u := url.URL{Scheme: "http", Host: r.ServerAddr, Path: r.UpdateURL}
	if r.CryptoKey != nil {
		payload, err = encryptData(payload, r.CryptoKey)
		if err != nil {
			return err
		}
	}
	if req, err = r.createRequest(u, payload); err != nil {
		return err
	}
	if err = r.addHMACHeader(req, payload); err != nil {
		return err
	}

	f := func() (any, error) {
		client := &http.Client{}
		var resp *http.Response
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, ErrResponseNotOK
		}
		defer resp.Body.Close()
		var body []byte
		body, err = io.ReadAll(resp.Body)
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

func (r *HTTPReportAdapter) createRequest(path url.URL, payload []byte) (*http.Request, error) {
	var compressed bytes.Buffer
	compressed, err := r.compress(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", path.String(), &compressed)
	if err != nil {
		r.Logger.Errorf("Failed to create a request: %v\n", err.Error())
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")
	if addr, err := getLocalIP(); err != nil {
		r.Logger.Errorf("Failed to add X-Real-IP header: %v", err)
		return nil, err
	} else {
		req.Header.Add("X-Real-IP", addr)
	}
	return req, nil
}

func (r *HTTPReportAdapter) addHMACHeader(req *http.Request, payload []byte) error {
	if hash, err := utils.HashData(payload, r.HMACKey); err != nil {
		return err
	} else if hash != "" {
		req.Header.Add("HashSHA256", hash)
	}
	return nil
}

func (r *HTTPReportAdapter) compress(payload []byte) (bytes.Buffer, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write(payload)
	if err != nil {
		return buf, err
	}
	if err := gz.Close(); err != nil {
		return buf, err
	}
	return buf, nil
}
