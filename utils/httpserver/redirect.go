package httpserver

import (
	"context"
	"net/http"

	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
)

func HandleRedirect(ctx context.Context, w http.ResponseWriter, status int, name string, location string, params any) {
	w.Header().Set(constants.HeaderLocation, location)
	switch responseFormat(w) {
	case m.ContentTypeJSON:
		// JSON clients get Location header only
		return
	case m.ContentTypeHTML:
		HandleResponse(ctx, w, status, name, params)
	}
}
