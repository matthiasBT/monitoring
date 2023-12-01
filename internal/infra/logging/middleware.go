// Package logging provides utilities for logging HTTP responses in Go.
// It includes a middleware function that can be used with the http package
// to log details about each HTTP request and its response.
package logging

import (
	"net/http"
	"time"
)

// responseMetadata holds the status code and size of an HTTP response.
type responseMetadata struct {
	status int // HTTP status code
	size   int // Size of the response body in bytes
}

// extendedWriter is an implementation of http.ResponseWriter that captures
// response metadata (status code and body size) in addition to providing
// standard response writing capabilities.
type extendedWriter struct {
	http.ResponseWriter                   // Embedded ResponseWriter to retain its methods
	response            *responseMetadata // Pointer to responseMetadata to store response details
}

// Write writes the data to the connection as part of an HTTP reply.
// It overrides the http.ResponseWriter's Write method to capture
// the size of the written data.
func (w *extendedWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.response.size += size
	return size, err
}

// WriteHeader sends an HTTP response header with the provided
// status code. It overrides the http.ResponseWriter's WriteHeader
// method to capture the status code.
func (w *extendedWriter) WriteHeader(status int) {
	w.response.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Middleware returns a middleware function for logging HTTP requests
// and responses. It wraps the provided http.Handler with logging
// functionalities using the provided logger.
func Middleware(logger ILogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		logFn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			extWriter := &extendedWriter{
				ResponseWriter: w,
				response: &responseMetadata{
					status: http.StatusOK,
					size:   0,
				},
			}
			next.ServeHTTP(extWriter, r)
			duration := time.Since(start)
			logger.Infof("Served: %s %s, %v\n", r.Method, r.RequestURI, duration)
			logger.Infof(
				"Response: [%d] %d bytes \n", extWriter.response.status, extWriter.response.size,
			)
		}
		return http.HandlerFunc(logFn)
	}
}
