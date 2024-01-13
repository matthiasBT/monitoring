package adapters

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"testing"

	"github.com/matthiasBT/monitoring/internal/infra/compression"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/secure"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
	"github.com/stretchr/testify/assert"
)

func TestHTTPReportAdapter_compress(t *testing.T) {
	reportAdapter := NewHTTPReportAdapter(
		logging.SetupLogger(),
		"0.0.0.0:8000",
		"/updates/",
		utils.Retrier{
			Attempts:         1,
			IntervalFirst:    0,
			IntervalIncrease: 0,
			Logger:           logging.SetupLogger(),
		},
		[]byte{},
		nil,
		0,
	)
	result, err := reportAdapter.compress([]byte("foobar"))
	if err != nil {
		t.Errorf("Failed to compress data: %v", err)
	}

	decompressed, err := compression.Decompress(io.NopCloser(bytes.NewReader(result.Bytes())))
	if err != nil {
		t.Errorf("Failed to decompress data: %v", decompressed)
	}
	assert.Equal(t, bytes.Equal(decompressed, []byte("foobar")), true)
}

func Test_encryptData(t *testing.T) {
	privateKey, err := generateRSA()
	if err != nil {
		t.Fatalf("Failed to create RSA key")
	}
	reportAdapter := NewHTTPReportAdapter(
		logging.SetupLogger(),
		"0.0.0.0:8000",
		"/updates/",
		utils.Retrier{
			Attempts:         1,
			IntervalFirst:    0,
			IntervalIncrease: 0,
			Logger:           logging.SetupLogger(),
		},
		[]byte{},
		&privateKey.PublicKey,
		0,
	)
	cipher, err := encryptData([]byte("foobar"), reportAdapter.CryptoKey)
	if err != nil {
		t.Errorf("Failed to encrypt data: %v", err)
	}
	decrypted, err := secure.Decrypt(cipher, privateKey)
	if err != nil {
		t.Errorf("Failed to decrypt data: %v", err)
	}
	assert.Equal(t, bytes.Equal(decrypted, []byte("foobar")), true)
}

func generateRSA() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
