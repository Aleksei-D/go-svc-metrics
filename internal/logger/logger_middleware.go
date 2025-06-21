package logger

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &responseData{}
		lw := &loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		next.ServeHTTP(lw, r)
		duration := time.Since(start)
		Log.Info("REQUEST", zap.String("METHOD", r.Method), zap.String("URI", r.URL.String()), zap.Duration("DURATION", duration))
		Log.Info("RESPONSE", zap.Int("STATUS", responseData.status), zap.Int("SIZE", responseData.size))
	})
}
