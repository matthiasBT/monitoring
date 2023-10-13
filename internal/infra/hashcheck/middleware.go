package hashcheck

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

type responseMetadata struct {
	data []byte
}

type extendedWriter struct {
	http.ResponseWriter
	response *responseMetadata
	hmacKey  string
}

func (w *extendedWriter) Write(b []byte) (int, error) {
	w.response.data = append(w.response.data, b...)
	serverHash, err := hashData([]byte(w.hmacKey), &w.response.data) // "{"id":"SD11","type":"counter","delta":1}"
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return len(err.Error()), err
	}
	w.Header().Set("HashSHA256", serverHash)
	size, err := w.ResponseWriter.Write(b)
	return size, err
}

func MiddlewareReader(key string) func(next http.Handler) http.Handler {
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

func MiddlewareWriter(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		addHashFn := func(w http.ResponseWriter, r *http.Request) {
			extWriter := &extendedWriter{
				ResponseWriter: w,
				response: &responseMetadata{
					data: []byte{},
				},
				hmacKey: key,
			}
			next.ServeHTTP(extWriter, r)
		}
		return http.HandlerFunc(addHashFn)
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
