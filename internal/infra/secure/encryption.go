// Package secure provides middleware for decryption of encrypted messages
package secure

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
)

// MiddlewareCryptoReader returns a middleware function that decrypts request body using an RSA private key
func MiddlewareCryptoReader(key *rsa.PrivateKey) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		checkHashFn := func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			decrypted, err := decrypt(body, key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(decrypted))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(checkHashFn)
	}
}

func decrypt(payload []byte, key *rsa.PrivateKey) ([]byte, error) {
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
