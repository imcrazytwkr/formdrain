package form

import (
	"net/http"
	"strings"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/httpserver"
	"github.com/imcrazytwkr/formdrain/httpserver/origin"
	"github.com/rs/zerolog"
)

func (r *formRouter) CheckCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log := zerolog.Ctx(req.Context())

		// Necessary for caching purposes
		w.Header().Set("Vary", constants.HeaderOrigin)

		originHost, err := origin.ParseOriginHost(req)
		if err != nil || len(originHost) == 0 {
			if err != nil {
				// Handling invalid Referer header as warning rather than full-blown error
				log.Warn().Err(err).Msg("Failed to parse CORS origin")
			}

			httpserver.HandleError(req.Context(), w, http.StatusForbidden, errInvalidOrigin)
			return
		}

		siteConfig, ok := getSiteConfig(req.Context())
		if !ok {
			log.Error().Msg("site config does not exist in context after being attached")
			httpserver.HandleStatus(req.Context(), w, http.StatusInternalServerError)
			return
		}

		if !strings.EqualFold(originHost, siteConfig.Hostname) {
			log.Warn().Str("request_origin", originHost).Str("expected_origin", siteConfig.Hostname).Msg("CORS origin mismatch")
			httpserver.HandleError(req.Context(), w, http.StatusForbidden, errInvalidOrigin)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", req.Header.Get(constants.HeaderOrigin))
		next.ServeHTTP(w, req)
	})
}

func (r *formRouter) HandlePreflight(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", http.MethodPost)
	w.Header().Set("Access-Control-Allow-Headers", constants.HeaderAccept+", "+constants.HeaderContentType)

	w.WriteHeader(http.StatusNoContent)
}
