package auth

import (
	"net/http"
	"time"

	"github.com/imcrazytwkr/formdrain/constants"
)

func (r *authRouter) sessionCookie(value string, expires time.Time, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     constants.CookieSession,
		Value:    value,
		Path:     "/",
		Domain:   r.config.CookieDomain,
		HttpOnly: true,
		Secure:   r.config.CookieSecure,
		SameSite: r.config.CookieSameSite.SameSite(),
		Expires:  expires,
		MaxAge:   maxAge,
	}
}
