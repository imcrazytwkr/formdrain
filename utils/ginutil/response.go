package ginutil

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
)

func HandleResponse(c *gin.Context, status int, name string, params any) {
	format := getResponseFormat(c)
	str, ok := params.(string)
	if ok {
		c.Header(constants.HeaderContentType, format.String())
		c.Status(status)
		io.WriteString(c.Writer, str)
		return
	}

	switch format {
	case m.ContentTypeJSON:
		c.JSON(status, params)
	case m.ContentTypeHTML:
		c.HTML(status, name, params)
	}
}
