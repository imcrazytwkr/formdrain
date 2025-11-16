package ginutil

import (
	"net/http"

	"github.com/gin-gonic/gin"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/rs/zerolog"
)

func HandleStatus(c *gin.Context, status int) {
	handleErrorMessage(c, status, "")
}

func HandleError(c *gin.Context, status int, err error) {
	if status < http.StatusInternalServerError {
		// Errors in 4xx range are user-generated so we proxy the error
		handleErrorMessage(c, status, err.Error())
		return
	}

	// Server errors should be replaced with their generic forms
	HandleStatus(c, status)
}

func handleErrorMessage(c *gin.Context, status int, message string) {
	log := zerolog.Ctx(c.Request.Context()).With().Str("handler", "error").Logger()

	if len(message) < 1 {
		message = http.StatusText(status)
	}

	if len(message) < 1 {
		// Mapping non-standard response status
		log.Warn().Int("status", status).Msg("unexpected status code")
		status = http.StatusInternalServerError
		message = http.StatusText(status)
	}

	switch getResponseFormat(c) {
	case m.ContentTypeJSON:
		c.JSON(status, &gin.H{
			"status":  status,
			"message": message,
		})
	case m.ContentTypeHTML:
		c.HTML(status, "errors/generic.html", &gin.H{
			"status":  status,
			"title":   http.StatusText(status),
			"message": message,
		})
	}

	c.Abort()
}
