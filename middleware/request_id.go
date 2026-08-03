package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/httpserver/clientip"
	"github.com/rs/zerolog"
)

func RequestId() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := ""

			// @NOTE: if this condition holds, we are running behind a trusted proxy
			if clientip.ClientIP(r) != clientip.RemoteIP(r) {
				requestId = r.Header.Get(constants.HeaderRequestID)
			}

			if len(requestId) < 1 {
				requestId = uuid.NewString()
			}

			log := zerolog.Ctx(r.Context()).With().Str("request_id", requestId).Logger()
			r = r.WithContext(log.WithContext(r.Context()))

			next.ServeHTTP(w, r)
		})
	}
}
