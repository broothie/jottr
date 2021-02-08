package logger

import (
	"fmt"
	"net/http"
	"time"
)

func Middleware(log *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestSize := r.ContentLength
			if requestSize < 0 {
				requestSize = 0
			}

			query := ""
			if len(r.URL.RawQuery) > 0 {
				query = fmt.Sprintf("?%s", r.URL.RawQuery)
			}

			// Make request
			recorder := newLoggerRecorder(w)
			before := time.Now()
			next.ServeHTTP(recorder, r)
			elapsed := time.Since(before)

			// Log after
			log.Info(
				fmt.Sprintf("%s %s%s %dB | %d %s %dB | %v",
					// Request
					r.Method,
					r.URL.Path,
					query,
					requestSize,
					//Response
					recorder.status,
					http.StatusText(recorder.status),
					recorder.bodyLength,
					// Timing
					elapsed,
				),
			)
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter
	status     int
	bodyLength int
}

func newLoggerRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{ResponseWriter: w, status: http.StatusOK}
}

func (lr *responseRecorder) WriteHeader(status int) {
	lr.status = status
	lr.ResponseWriter.WriteHeader(status)
}

func (lr *responseRecorder) Write(body []byte) (int, error) {
	lr.bodyLength = len(body)
	return lr.ResponseWriter.Write(body)
}
