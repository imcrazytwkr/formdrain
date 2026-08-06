package httpserver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/rs/zerolog"
)

func HandleResponse(ctx context.Context, w http.ResponseWriter, status int, name string, params any) {
	format := ResponseFormat(w)

	str, ok := params.(string)
	if ok {
		w.Header().Set(constants.HeaderContentType, format.String())
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
	w.Header().Set(constants.HeaderContentType, m.ContentTypeJSON.String())
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

type templateDataProvider interface {
	TemplateData() map[string]any
}

func writeHTML(ctx context.Context, w http.ResponseWriter, status int, name string, params any) {
	w.Header().Set(constants.HeaderContentType, m.ContentTypeHTML.String())
	w.WriteHeader(status)

	log := zerolog.Ctx(ctx)

	tmpl, ok := templates[name]
	if !ok || tmpl == nil {
		log.Error().Str("template", name).Msg("template not found")
		_, err := io.WriteString(w, http.StatusText(status))
		if err != nil {
			log.Error().Err(err).Msg("failed to write fallback plain-text response")
		}
		return
	}

	var data map[string]any
	switch p := params.(type) {
	case templateDataProvider:
		data = p.TemplateData()
	case map[string]any:
		data = p
	default:
		if params != nil {
			log.Warn().
				Str("template", name).
				Type("params_type", params).
				Msg("unsupported HTML template params; using empty data")
		}

		data = map[string]any{}
	}

	err := tmpl.FRender(w, data)
	if err == nil {
		return
	}

	log.Error().Err(err).Msg("failed to render HTML response")

	_, err = io.WriteString(w, http.StatusText(status))
	if err != nil {
		log.Error().Err(err).Msg("failed to write fallback plain-text response")
	}
}
