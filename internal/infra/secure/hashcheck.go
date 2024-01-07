// Package secure provides middleware for adding and verifying HMAC SHA256
// hashes in HTTP requests and responses. It ensures the integrity of the data
// transmitted over HTTP.
package secure

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

// responseMetadata holds the data of an HTTP response to be hashed.
type responseMetadata struct {
	// data is the byte slice that stores the response data.
	data []byte
}

// extendedWriter is a wrapper around http.ResponseWriter that adds HMAC SHA256
// hash to the response header after writing the response body.
type extendedWriter struct {
	// ResponseWriter is embedded and allows extendedWriter to implement
	// the http.ResponseWriter interface.
	http.ResponseWriter

	// response holds the metadata of the response being written.
	response *responseMetadata

	// hmacKey is the secret key used for generating the HMAC hash.
	hmacKey string
}

// Write hashes the response data using HMAC SHA256 and writes it to the client.
// It also sets the HashSHA256 header in the response.
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

// MiddlewareHashReader returns a middleware function that verifies the HMAC SHA256
// hash of the request body. It compares the client-provided hash in the header
// with the server-generated hash to ensure data integrity.
func MiddlewareHashReader(key string) func(next http.Handler) http.Handler {
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

// MiddlewareHashWriter returns a middleware function that adds an HMAC SHA256
// hash to the response header. It uses extendedWriter to automatically hash
// the response data and append the hash to the response headers.
func MiddlewareHashWriter(key string) func(next http.Handler) http.Handler {
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
