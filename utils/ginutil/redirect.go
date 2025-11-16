package ginutil

import (
	"github.com/gin-gonic/gin"
	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
)

func HandleRedirect(c *gin.Context, status int, name string, location string, params any) {
	c.Header(constants.HeaderLocation, location)
	switch getResponseFormat(c) {
	case m.ContentTypeJSON:
		break
	case m.ContentTypeHTML:
		HandleResponse(c, status, name, params)
	}

}
