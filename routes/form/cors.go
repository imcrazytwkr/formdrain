package form

import (
	"net/http"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/utils/httpserver"
	"github.com/rs/zerolog"
)

func (r *formRouter) CheckCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log := zerolog.Ctx(req.Context())

		originHost, err := httpserver.ParseOriginHost(req)
		if len(originHost) < 1 {
			if err != nil {
				httpserver.HandleError(req.Context(), w, http.StatusForbidden, err)
			}
			return
		}

		if err != nil {
			// Handling invalid Referer header as warning rather than full-blown error
			log.Warn().Msg(err.Error())
		}

		siteConfig, ok := getSiteConfig(req.Context())
		if !ok {
			log.Error().Msg("site config does not exist in context after being attached")
			httpserver.HandleStatus(req.Context(), w, http.StatusInternalServerError)
			return
		}

		if originHost != siteConfig.Hostname {
			log.Warn().Str("request_origin", originHost).Str("expected_origin", siteConfig.Hostname).Msg("CORS origin mismatch")
			httpserver.HandleError(req.Context(), w, http.StatusForbidden, errInvalidOrigin)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", originHost)

		// First call is Set to clear previous values, if any
		w.Header().Set("Vary", constants.HeaderOrigin)
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")

		next.ServeHTTP(w, req)
	})
}

func (r *formRouter) HandlePreflight(w http.ResponseWriter, req *http.Request) {
	// First call is Set to clear previous values, if any
	w.Header().Set("Access-Control-Allow-Methods", http.MethodOptions)
	w.Header().Add("Access-Control-Allow-Methods", http.MethodPost)

	w.WriteHeader(http.StatusNoContent)
}
