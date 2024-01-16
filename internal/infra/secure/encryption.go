// Package secure provides middleware for decryption of encrypted messages
package secure

import (
	"bytes"
	"crypto/rsa"
	"io"
	"net/http"

	"github.com/matthiasBT/monitoring/internal/infra/utils"
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
			decrypted, err := utils.Decrypt(body, key)
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
