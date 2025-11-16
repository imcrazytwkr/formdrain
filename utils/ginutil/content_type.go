package ginutil

import (
	"github.com/gin-gonic/gin"
	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
)

func GetContentType(c *gin.Context) m.ContentType {
	raw, ok := c.Get(constants.KeyContentType)
	if !ok {
		return m.ContentTypeUndefined
	}

	contentType, ok := raw.(m.ContentType)
	if !ok {
		return m.ContentTypeUndefined
	}

	return contentType
}
