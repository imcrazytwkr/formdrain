package form

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imcrazytwkr/formdrain/utils/ginutil"
	"github.com/rs/zerolog"
)

func (r *formRouter) CheckCORS(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	originHost, err := ginutil.ParseOriginHost(c)
	if len(originHost) < 1 {
		if err != nil {
			ginutil.HandleError(c, http.StatusForbidden, err)
		}

		return
	}

	if err != nil {
		// Handling invalid Referer header as warning rather than full-blown error
		log.Warn().Msg(err.Error())
	}

	siteConfig, ok := getSiteConfig(c)
	if !ok {
		log.Error().Msg("site config does not exist in context after being attached")
		ginutil.HandleStatus(c, http.StatusInternalServerError)
		return
	}

	if originHost != siteConfig.Hostname {
		log.Warn().Str("request_origin", originHost).Str("expected_origin", siteConfig.Hostname).Msg("CORS origin mismatch")
		ginutil.HandleError(c, http.StatusForbidden, errInvalidOrigin)
		return
	}

	c.Header("Access-Control-Allow-Origin", originHost)

	// First call is Set to clear previous values, if any
	c.Writer.Header().Set("Vary", "Origin")
	c.Writer.Header().Add("Vary", "Access-Control-Request-Method")
	c.Writer.Header().Add("Vary", "Access-Control-Request-Headers")
}

func (r *formRouter) HandlePreflight(c *gin.Context) {
	// First call is Set to clear previous values, if any
	c.Writer.Header().Set("Access-Control-Allow-Methods", http.MethodOptions)
	c.Writer.Header().Add("Access-Control-Allow-Methods", http.MethodPost)

	// Empty response body on success
	c.AbortWithStatus(http.StatusNoContent)
}
