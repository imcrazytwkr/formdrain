package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/httpserver"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/models/session"
	"github.com/imcrazytwkr/formdrain/routes/auth/mappers"
	"github.com/imcrazytwkr/formdrain/services/account"
	"github.com/imcrazytwkr/formdrain/utils/bodyparser"
	"github.com/imcrazytwkr/formdrain/utils/maputil"
)

const sessionTTL = 24 * time.Hour

func (r *authRouter) HandleLogin(w http.ResponseWriter, req *http.Request) {
	log := getLoggerForAction(req.Context(), actionLogin)
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

	a, err := r.accountService.Login(ctx, maputil.String(data, "email"), maputil.String(data, "password"))
	if err != nil {
		if errors.Is(err, account.ErrInvalidCredentials) {
			httpserver.HandleError(ctx, w, http.StatusUnauthorized, err)
			return
		}

		log.Err(err).Msg("authorisation error")
		httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
		return
	}

	session := &session.Session{
		AccountID: a.ID,
		ExpiresAt: time.Now().Add(sessionTTL),
	}

	err = r.sessionRepository.Create(ctx, session)
	if err != nil {
		httpserver.HandleError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     constants.CookieSession,
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  session.ExpiresAt,
		MaxAge:   max(int(time.Until(session.ExpiresAt).Seconds()), 1),
	})

	if httpserver.ResponseFormat(w) == m.ContentTypeJSON {
		httpserver.HandleResponse(ctx, w, http.StatusCreated, "", mappers.User(a))
		return
	}

	origin, _ := maputil.GetString(data, "origin")
	httpserver.HandleRedirect(
		ctx,
		w,
		http.StatusSeeOther,
		"form/redirect",
		origin,
		map[string]any{"redirect_to": SanitizeAbsolutePath(origin)},
	)
}
