package middleware

import (
	"net/http"

	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/utils/httpserver"
)

func ResponseFormatParser(allowed ...m.ContentType) func(http.Handler) http.Handler {
	offers := make([]m.ContentType, 0, len(allowed))
	for _, v := range allowed {
		if v != m.ContentTypeUndefined {
			offers = append(offers, v)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			format := httpserver.Negotiate(r, offers)
			if format != m.ContentTypeUndefined {
				w.Header().Set(constants.HeaderContentType, format.String())
			}

			next.ServeHTTP(w, r)
		})
	}
}
