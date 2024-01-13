// Package secure provides middleware for filtering clients' requests based on their IP addresses
package secure

import (
	"net"
	"net/http"

	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
)

// MiddlewareIPFilter either allows a client to proceed with its request or blocks it if the client's IP is not trusted
func MiddlewareIPFilter(logger logging.ILogger, rawSubnet string) func(next http.Handler) http.Handler {
	subnet := utils.ParseSubnet(rawSubnet)
	return func(next http.Handler) http.Handler {
		checkIP := func(w http.ResponseWriter, r *http.Request) {
			clientIPRaw := r.Header.Get("X-Real-IP")
			logger.Infof("Client IP address: %s", clientIPRaw)
			if clientIPRaw == "" {
				http.Error(w, "Missing X-Real-IP header value", http.StatusForbidden)
				return
			}
			var clientIP = net.ParseIP(clientIPRaw)
			if clientIP == nil {
				http.Error(w, "Invalid X-Real-IP header value", http.StatusForbidden)
				return
			}
			if !subnet.Contains(clientIP) {
				http.Error(w, "Untrusted client IP address", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(checkIP)
	}
}
