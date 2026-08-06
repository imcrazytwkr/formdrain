package middleware

import (
	"context"
	"net/http"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/httpserver"
	"github.com/imcrazytwkr/formdrain/models/session"
	"github.com/imcrazytwkr/formdrain/repositories"
)

type sessionContextKey struct{}

func Authenticated(sessions repositories.SessionRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(constants.CookieSession)
			if err != nil || len(cookie.Value) < 1 {
				httpserver.HandleStatus(r.Context(), w, http.StatusUnauthorized)
				return
			}

			sess, err := sessions.GetByID(r.Context(), cookie.Value)
			if err != nil {
				httpserver.HandleStatus(r.Context(), w, http.StatusInternalServerError)
				return
			}

			if sess == nil {
				httpserver.HandleStatus(r.Context(), w, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), sessionContextKey{}, sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SessionFromContext(ctx context.Context) (*session.Session, bool) {
	s, ok := ctx.Value(sessionContextKey{}).(*session.Session)
	return s, ok && s != nil
}
