package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

func Decrypt(payload []byte, key *rsa.PrivateKey) ([]byte, error) {
	encryptedKey, encryptedData := payload[:256], payload[256:] // RSA-2048
	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, encryptedKey, nil)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	if len(encryptedData) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := encryptedData[:aes.BlockSize]
	encryptedData = encryptedData[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(encryptedData))
	stream.XORKeyStream(plaintext, encryptedData)

	return plaintext, nil
}
