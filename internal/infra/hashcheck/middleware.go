package hashcheck

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

func Middleware(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		checkHashFn := func(w http.ResponseWriter, r *http.Request) {
			var clientHash string
			if clientHash = r.Header.Get("HashSHA256"); clientHash == "" {
				next.ServeHTTP(w, r)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			serverHash, err := hashData([]byte(key), &body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if clientHash != serverHash {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(body))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(checkHashFn)
	}
}

func hashData(key []byte, payload *[]byte) (string, error) {
	mac := hmac.New(sha256.New, key)
	if _, err := mac.Write(*payload); err != nil {
		return "", err
	}
	hash := mac.Sum(nil)
	result := hex.EncodeToString(hash)
	return result, nil
}
