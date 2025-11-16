package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const requestIdHeader = "X-Request-ID"

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := ""

		// @NOTE: if this condition holds, we are running behind a trusted proxy
		if c.ClientIP() != c.RemoteIP() {
			requestId = c.GetHeader(requestIdHeader)
		}

		if len(requestId) < 1 {
			requestId = uuid.NewString()
		}

		log := zerolog.Ctx(c.Request.Context()).With().Str("request_id", requestId).Logger()
		c.Request = c.Request.WithContext(log.WithContext(c.Request.Context()))

		c.Next()
	}
}
