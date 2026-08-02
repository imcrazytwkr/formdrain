package httpserver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/templates"
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
		writeHTML(ctx, w, status, GetTemplate(name), params)
	}
}

func HandleResponseTemplate(ctx context.Context, w http.ResponseWriter, status int, template templates.Template, params any) {
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
		writeHTML(ctx, w, status, template, params)
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

type templateDataProvider interface {
	TemplateData() map[string]any
}

func writeHTML(ctx context.Context, w http.ResponseWriter, status int, template templates.Template, params any) {
	if template == nil {
		panic("httpserver: writeHTML called with nil template")
	}

	setResponseContentType(w, m.ContentTypeHTML)
	w.WriteHeader(status)

	log := zerolog.Ctx(ctx)

	var data map[string]any
	switch p := params.(type) {
	case templateDataProvider:
		data = p.TemplateData()
	case map[string]any:
		data = p
	default:
		if params != nil {
			log.Warn().
				// Str("template", name).
				Type("params_type", params).
				Msg("unsupported HTML template params; using empty data")
		}

		data = map[string]any{}
	}

	err := template.Execute(w, data)
	if err == nil {
		return
	}

	log.Error().Err(err).Msg("failed to render HTML response")

	_, err = io.WriteString(w, http.StatusText(status))
	if err != nil {
		log.Error().Err(err).Msg("failed to write fallback plain-text response")
	}
}
