package middleware

import (
	"fmt"
	"net/http"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/httpserver"
	"github.com/imcrazytwkr/formdrain/httpserver/contenttype"
	m "github.com/imcrazytwkr/formdrain/models/http"
)

func ContentTypeParser() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch:
				contentType := m.ParseFormContentType(contenttype.RequestContentType(r))
				if contentType == m.ContentTypeUndefined {
					httpserver.HandleError(
						r.Context(),
						w,
						http.StatusUnsupportedMediaType,
						fmt.Errorf("valid Content-Type header should be specified for %s requests", r.Method),
					)
					return
				}

				r = r.WithContext(contenttype.WithContentType(r.Context(), contentType))

				// Guess response format if Accept did not set Content-Type.
				if len(w.Header().Get(constants.HeaderContentType)) < 1 {
					w.Header().Set(constants.HeaderContentType, contentType.String())
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
