package ginutil

import (
	"github.com/gin-gonic/gin"
	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
)

func GetResponseFormat(c *gin.Context) m.ContentType {
	raw, ok := c.Get(constants.KeyResponseFormat)
	if !ok {
		return m.ContentTypeUndefined
	}

	contentType, ok := raw.(m.ContentType)
	if !ok {
		return m.ContentTypeUndefined
	}

	return contentType
}

// @NOTE: this is specific version for response/error handlers
func getResponseFormat(c *gin.Context) m.ContentType {
	format := GetResponseFormat(c)
	if format == m.ContentTypeUndefined {
		return m.ContentTypeHTML
	}

	return format
}
