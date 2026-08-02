package httpserver

import (
	"context"
	"net/http"

	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/templates"
)

func HandleRedirect(ctx context.Context, w http.ResponseWriter, status int, name string, location string, params any) {
	HandleRedirectTemplate(ctx, w, status, GetTemplate(name), location, params)
}

func HandleRedirectTemplate(ctx context.Context, w http.ResponseWriter, status int, template templates.Template, location string, params any) {
	w.Header().Set(constants.HeaderLocation, location)
	switch responseFormat(w) {
	case m.ContentTypeJSON:
		// JSON clients get Location header only
		return
	case m.ContentTypeHTML:
		HandleResponseTemplate(ctx, w, status, template, params)
	}
}
