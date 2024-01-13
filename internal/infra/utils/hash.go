package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HashData(payload, hmacKey []byte) (string, error) {
	if bytes.Equal(hmacKey, []byte{}) {
		return "", nil
	}
	mac := hmac.New(sha256.New, hmacKey)
	if _, err := mac.Write(payload); err != nil {
		return "", err
	}
	hash := mac.Sum(nil)
	result := hex.EncodeToString(hash)
	return result, nil
}
