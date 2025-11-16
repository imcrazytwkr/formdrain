package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/utils/ginutil"
)

func ContentTypeParser(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			contentType := m.ParseFormContentType(c.ContentType())
			if contentType == m.ContentTypeUndefined {
				// No Content-Type specified for mutable body
				ginutil.HandleError(
					c,
					http.StatusUnsupportedMediaType,
					fmt.Errorf("valid Content-Type header should be specified for %s requests", c.Request.Method),
				)

				c.Abort()
				return
			}

			c.Set(constants.KeyContentType, contentType)

			// Guessing response format if no Accepts header was supplied based on
			// Content-Type header of request
			format := ginutil.GetResponseFormat(c)
			if format == m.ContentTypeUndefined {
				c.Set(constants.KeyResponseFormat, contentType)
			}

			fallthrough
		default:
			c.Next()
		}
	}
}
