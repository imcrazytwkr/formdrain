package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func DefaultLogger() gin.HandlerFunc {
	return LoggerWithConfig(&log.Logger, nil)
}

// LoggerWithConfig instance a Logger middleware with config.
func LoggerWithConfig(logger *zerolog.Logger, skipPaths []string) gin.HandlerFunc {
	var skip map[string]struct{}

	skipped := len(skipPaths)
	if skipped > 0 {
		skip = make(map[string]struct{}, skipped)

		for _, path := range skipPaths {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		logger.With()
		c.Request = c.Request.WithContext(logger.WithContext(c.Request.Context()))

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log only when path is not being skipped
		_, ok := skip[path]
		if ok {
			return
		}

		stop := time.Now()
		statusCode := c.Writer.Status()

		if len(raw) > 0 {
			path = path + "?" + raw
		}

		// Log using the params
		var logEvent *zerolog.Event
		if statusCode < http.StatusInternalServerError {
			logEvent = logger.Info()
		} else {
			logEvent = logger.Error()
		}

		logEvent.Str("client_ip", c.ClientIP())
		logEvent.Str("method", c.Request.Method)
		logEvent.Int("status_code", statusCode)
		logEvent.Int("body_size", c.Writer.Size())
		logEvent.Str("path", path)
		logEvent.Str("latency", stop.Sub(start).String())
		logEvent.Msg(c.Errors.ByType(gin.ErrorTypePrivate).String())
	}
}
