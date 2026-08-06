package auth

import (
	"errors"
	"net/http"

	"github.com/imcrazytwkr/formdrain/httpserver"
	m "github.com/imcrazytwkr/formdrain/models/http"
)

func (r *authRouter) HandleLoginForm(w http.ResponseWriter, req *http.Request) {
	log := getLoggerForAction(req.Context(), actionForm)
	ctx := log.WithContext(req.Context())

	if httpserver.ResponseFormat(w) == m.ContentTypeJSON {
		httpserver.HandleError(ctx, w, http.StatusNotAcceptable, errors.New("login form is only available as HTML"))
		return
	}

	origin := SanitizeAbsolutePath(req.URL.Query().Get("redirect"))
	httpserver.HandleResponse(ctx, w, http.StatusOK, "auth/login", map[string]any{
		"origin": origin,
	})
}
