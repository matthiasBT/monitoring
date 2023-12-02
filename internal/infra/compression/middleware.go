// Package compression provides middleware for handling gzip compression
// and decompression in HTTP requests and responses.
package compression

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// gzipWriter is a wrapper around http.ResponseWriter that writes
// response data using gzip compression.
type gzipWriter struct {
	// ResponseWriter is embedded and allows gzipWriter to implement
	// the http.ResponseWriter interface.
	http.ResponseWriter

	// Writer is the gzip writer used to compress the response data.
	Writer io.Writer
}

// Write compresses the given bytes using gzip and writes them to the response.
// It implements the Write method of the http.ResponseWriter interface.
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// MiddlewareWriter is a middleware function that handles gzip compression
// for HTTP responses. If the client accepts gzip encoding, it compresses
// the response, otherwise it passes the response through unchanged.
func MiddlewareWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo: split string and check
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

// MiddlewareReader is a middleware function that handles gzip decompression
// for HTTP requests. If the request is gzip-encoded, it decompresses the
// request body, otherwise it passes the request through unchanged.
func MiddlewareReader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo: split string and check
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		body, err := io.ReadAll(gz)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	})
}
