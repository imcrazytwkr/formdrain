package apiv1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/middleware"
	models "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/repositories"
	"github.com/imcrazytwkr/formdrain/routes"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/api"
)

type apiV1Router struct {
	sessions repositories.SessionRepository
	accounts repositories.AccountRepository
	sites    repositories.SiteConfigRepository
	forms    repositories.FormConfigRepository
}

func NewApiV1Router(
	sessions repositories.SessionRepository,
	accounts repositories.AccountRepository,
	sites repositories.SiteConfigRepository,
	forms repositories.FormConfigRepository,
) routes.RouteContainer {
	return &apiV1Router{
		sessions: sessions,
		accounts: accounts,
		sites:    sites,
		forms:    forms,
	}
}

func (r *apiV1Router) Router(router chi.Router) {
	router.Use(forceJSON)
	router.Use(middleware.Authenticated(r.sessions))
	api.HandlerFromMux(api.NewStrictHandler(r, nil), router)
}

func forceJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set(constants.HeaderContentType, models.ContentTypeJSON.String())
		next.ServeHTTP(w, req)
	})
}
