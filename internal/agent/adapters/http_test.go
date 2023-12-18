package adapters

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/secure"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
	"github.com/stretchr/testify/assert"
)

func Test_encryptData(t *testing.T) {
	hmacKey := make([]byte, 32)
	if _, err := rand.Read(hmacKey); err != nil {
		t.Fatalf("Failed to create HMAC key: %v", err)
	}
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
		hmacKey,
		&privateKey.PublicKey,
		0,
	)
	reportAdapter.Jobs = make(chan []byte, 1)
	cipher, err := reportAdapter.encryptData([]byte("foobar"))
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
