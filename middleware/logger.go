package middleware

import (
	"log"
	"net/http"
	"time"
)

var requestLogger *log.Logger

// SetLogger sets the logger used by the Logger middleware.
func SetLogger(logger *log.Logger) {
	requestLogger = logger
}

// Logger is a middleware that logs the request method, response status, URL path, and the duration of the request.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestLogger == nil {
			requestLogger = log.Default()
		}
		start := time.Now()
		prw := newProxyWriter(w)
		next.ServeHTTP(prw, r)
		requestLogger.Printf("%s: %d %s %d ms", r.Method, prw.Status(), r.URL.EscapedPath(), time.Since(start).Microseconds())
	})
}
