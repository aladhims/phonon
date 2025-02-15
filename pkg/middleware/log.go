package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// LoggingMiddleware logs incoming HTTP requests with timing information
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		srw := &statusResponseWriter{ResponseWriter: w}

		next.ServeHTTP(srw, r)

		logrus.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     srw.status,
			"duration":   time.Since(start),
			"user_agent": r.UserAgent(),
			"remote_ip":  r.RemoteAddr,
		}).Info("Handled request")
	})
}
