package httpserver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/rs/zerolog"
)

func HandleResponse(ctx context.Context, w http.ResponseWriter, status int, name string, params any) {
	format := responseFormat(w)

	str, ok := params.(string)
	if ok {
		setResponseContentType(w, format)
		w.WriteHeader(status)
		io.WriteString(w, str)
		return
	}

	switch format {
	case m.ContentTypeJSON:
		writeJSON(ctx, w, status, params)
	case m.ContentTypeHTML:
		writeHTML(ctx, w, status, name, params)
	}
}

func writeJSON(ctx context.Context, w http.ResponseWriter, status int, params any) {
	setResponseContentType(w, m.ContentTypeJSON)
	w.WriteHeader(status)

	body, err := json.Marshal(params)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to marshal JSON response")
		return
	}

	_, err = w.Write(body)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to write JSON response")
	}
}

func writeHTML(ctx context.Context, w http.ResponseWriter, status int, name string, params any) {
	setResponseContentType(w, m.ContentTypeHTML)
	w.WriteHeader(status)

	log := zerolog.Ctx(ctx)
	if templates == nil {
		log.Error().Str("template", name).Msg("templates are not loaded")
		_, err := io.WriteString(w, http.StatusText(status))
		if err != nil {
			log.Error().Err(err).Msg("failed to write fallback plain-text response")
		}
		return
	}

	err := templates.ExecuteTemplate(w, name, params)
	if err == nil {
		return
	}

	log.Error().Err(err).Msg("failed to render HTML response")

	_, err = io.WriteString(w, http.StatusText(status))
	if err != nil {
		log.Error().Err(err).Msg("failed to write fallback plain-text response")
	}
}
