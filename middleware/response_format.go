package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
)

func ResponseFormatParser(allowed ...m.ContentType) gin.HandlerFunc {
	textContentTypes := make([]string, len(allowed)+1)
	textContentTypes[0] = m.ContentTypeUndefined.String()
	i := 1

	for _, v := range allowed {
		if v != m.ContentTypeUndefined {
			textContentTypes[i] = v.String()
			i++
		}
	}

	if i < len(allowed) {
		textContentTypes = textContentTypes[0:i]
	}

	return func(c *gin.Context) {
		c.Set(constants.KeyResponseFormat, m.ParseContentType(c.NegotiateFormat(textContentTypes...)))
		c.Next()
	}
}
