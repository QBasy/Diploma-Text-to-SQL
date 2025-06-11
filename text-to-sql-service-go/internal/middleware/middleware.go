package middleware

import (
	"net/http"
	"text-to-sql/pkg/logger"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: 200}
		start := time.Now()

		next.ServeHTTP(rec, r)

		duration := time.Since(start)
		logger.InfoLogger.Printf("[%s] %s %s %d %s", r.Method, r.RequestURI, r.RemoteAddr, rec.status, duration)
	})
}
