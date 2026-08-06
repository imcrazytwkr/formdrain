package httpserver

import (
	"context"
	"net/http"

	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/rs/zerolog"
)

func HandleStatus(ctx context.Context, w http.ResponseWriter, status int) {
	handleErrorMessage(ctx, w, status, "")
}

func HandleError(ctx context.Context, w http.ResponseWriter, status int, err error) {
	if status < http.StatusInternalServerError {
		// Errors in 4xx range are user-generated so we proxy the error
		handleErrorMessage(ctx, w, status, err.Error())
		return
	}

	// Server errors should be replaced with their generic forms
	HandleStatus(ctx, w, status)
}

func handleErrorMessage(ctx context.Context, w http.ResponseWriter, status int, message string) {
	log := zerolog.Ctx(ctx).With().Str("handler", "error").Logger()

	if len(message) < 1 {
		message = http.StatusText(status)
	}

	if len(message) < 1 {
		// Mapping non-standard response status
		log.Warn().Int("status", status).Msg("unexpected status code")
		status = http.StatusInternalServerError
		message = http.StatusText(status)
	}

	format := ResponseFormat(w)
	switch format {
	case m.ContentTypeJSON:
		writeJSON(ctx, w, status, map[string]any{
			"status":  status,
			"message": message,
		})
	case m.ContentTypeHTML:
		writeHTML(ctx, w, status, "errors/generic", map[string]any{
			"status":  status,
			"title":   http.StatusText(status),
			"message": message,
		})
	}
}
