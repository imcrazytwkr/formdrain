package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/httpserver"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/utils/bodyparser"
	"github.com/imcrazytwkr/formdrain/utils/maputil"
)

var timeEpochStart = time.Unix(0, 0).UTC()

func (r *authRouter) HandleLogout(w http.ResponseWriter, req *http.Request) {
	log := getLoggerForAction(req.Context(), actionLogout)
	ctx := log.WithContext(req.Context())

	data, err := bodyparser.Parse(req)
	if err != nil {
		if errors.Is(err, bodyparser.ErrBodyTooLarge) {
			httpserver.HandleError(ctx, w, http.StatusRequestEntityTooLarge, err)
			return
		}

		httpserver.HandleError(ctx, w, http.StatusBadRequest, err)
		return
	}

	cookie, err := req.Cookie(constants.CookieSession)
	if err == nil && len(cookie.Value) > 0 {
		err := r.sessionRepository.Delete(ctx, cookie.Value)
		if err != nil {
			log.Err(err).Msg("failed to delete session")
			httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
			return
		}
	}

	http.SetCookie(w, r.sessionCookie("", timeEpochStart, -1))

	if httpserver.ResponseFormat(w) == m.ContentTypeJSON {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	origin := SanitizeAbsolutePath(maputil.String(data, "origin"))
	httpserver.HandleRedirect(
		ctx,
		w,
		http.StatusSeeOther,
		"form/redirect",
		origin,
		map[string]any{"redirect_to": origin},
	)
}
