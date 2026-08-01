package middleware

import (
	"net/http"
	"time"

	"github.com/imcrazytwkr/formdrain/utils/httpserver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func DefaultLogger() func(http.Handler) http.Handler {
	return LoggerWithConfig(&log.Logger, nil)
}

// LoggerWithConfig returns a logger middleware with optional path skips.
func LoggerWithConfig(logger *zerolog.Logger, skipPaths []string) func(http.Handler) http.Handler {
	var skip map[string]struct{}

	skipped := len(skipPaths)
	if skipped > 0 {
		skip = make(map[string]struct{}, skipped)
		for _, path := range skipPaths {
			skip[path] = struct{}{}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			raw := r.URL.RawQuery

			r = r.WithContext(logger.WithContext(r.Context()))
			ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			start := time.Now()
			next.ServeHTTP(ww, r)

			if _, ok := skip[path]; ok {
				return
			}

			if len(raw) > 0 {
				path = path + "?" + raw
			}

			statusCode := ww.status
			var logEvent *zerolog.Event
			if statusCode < http.StatusInternalServerError {
				logEvent = logger.Info()
			} else {
				logEvent = logger.Error()
			}

			logEvent.
				Str("client_ip", httpserver.ClientIP(r)).
				Str("method", r.Method).
				Int("status_code", statusCode).
				Int("body_size", ww.bytes).
				Str("path", path).
				Str("latency", time.Since(start).String()).
				Msg("")
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}
