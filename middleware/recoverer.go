package middleware

import (
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/httpserver"
	"github.com/rs/zerolog"
)

func Recoverer() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				rvr := recover()
				switch rvr {
				case nil:
					return
				case http.ErrAbortHandler:
					panic(rvr)
				default:
					zerolog.Ctx(r.Context()).Error().
						Interface("panic", rvr).
						Bytes("stack", debug.Stack()).
						Msg("panic recovered")
				}

				if strings.EqualFold(r.Header.Get(constants.HeaderConnection), "Upgrade") {
					return
				}

				httpserver.HandleStatus(r.Context(), w, http.StatusInternalServerError)
			}()

			next.ServeHTTP(w, r)
		})
	}
}
