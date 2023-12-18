// Package adapters provides functionalities to handle HTTP communication,
// specifically for reporting metrics. It includes structures and methods
// for sending reports, handling retries, and ensuring data integrity.
package adapters

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
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

	// jobs is an internal channel used to queue payloads for reporting.
	jobs chan []byte

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

// ErrResponseNotOK is an error indicating that the HTTP response status is not OK (200).
var ErrResponseNotOK = errors.New("response not OK")

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
	jobs := make(chan []byte, workerNum)
	adapter := HTTPReportAdapter{
		Logger:     logger,
		ServerAddr: serverAddr,
		UpdateURL:  updateURL,
		Retrier:    retrier,
		HMACKey:    hmacKey,
		CryptoKey:  cryptoKey,
		jobs:       jobs,
	}
	var i uint
	for i = 0; i < workerNum; i++ {
		go func() {
			for {
				data := <-jobs
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
	r.jobs <- payload
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
	r.jobs <- payload
	return nil
}

func (r *HTTPReportAdapter) report(payload []byte) error {
	var (
		req *http.Request
		err error
	)
	u := url.URL{Scheme: "http", Host: r.ServerAddr, Path: r.UpdateURL}
	if r.CryptoKey != nil {
		payload, err = r.encryptData(payload)
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
			r.Logger.Errorf("Request failed: %v\n", err.Error())
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			r.Logger.Errorf("Request failed with code: %d\n", resp.StatusCode)
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
	return req, nil
}

func (r *HTTPReportAdapter) addHMACHeader(req *http.Request, payload []byte) error {
	if hash, err := r.hashData(payload); err != nil {
		return err
	} else if hash != "" {
		req.Header.Add("HashSHA256", hash)
	}
	return nil
}

func (r *HTTPReportAdapter) hashData(payload []byte) (string, error) {
	if bytes.Equal(r.HMACKey, []byte{}) {
		return "", nil
	}
	mac := hmac.New(sha256.New, r.HMACKey)
	if _, err := mac.Write(payload); err != nil {
		r.Logger.Errorf("Failed to calculate hash: %v", err.Error())
		return "", err
	}
	hash := mac.Sum(nil)
	result := hex.EncodeToString(hash)
	r.Logger.Infof("HMAC-SHA256 hash: %s\n", result)
	return result, nil
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

func (r *HTTPReportAdapter) encryptData(payload []byte) ([]byte, error) {
	key, encryptedPayload, err := encryptAES(payload)
	if err != nil {
		r.Logger.Errorf("Error encrypting message: %v", err)
		return nil, err
	}
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, r.CryptoKey, key, nil)
	if err != nil {
		r.Logger.Errorf("Error encrypting AES key: %v", err)
		return nil, err
	}
	return append(encryptedKey, encryptedPayload...), nil
}

func encryptAES(plaintext []byte) ([]byte, []byte, error) {
	key := make([]byte, 32) // AES-256
	if _, err := rand.Read(key); err != nil {
		return nil, nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return nil, nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return key, ciphertext, nil
}
