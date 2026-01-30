package middleware

import (
	"net/http"
	"time"

	"log/slog"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	written int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

// Logging logs each request: method, path, status, duration, request_id.
// Uses slog; level is Info for 2xx, Warn for 4xx, Error for 5xx.
func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(wrapped, r)
			dur := time.Since(start)
			reqID := ""
			if id := r.Context().Value(RequestIDKey); id != nil {
				reqID, _ = id.(string)
			}
			attrs := []any{
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", wrapped.status),
				slog.Duration("duration", dur),
				slog.String("request_id", reqID),
			}
			if wrapped.status >= 500 {
				logger.Error("request", attrs...)
			} else if wrapped.status >= 400 {
				logger.Warn("request", attrs...)
			} else {
				logger.Info("request", attrs...)
			}
		})
	}
}
