package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/imcrazytwkr/formdrain/middleware"
	"github.com/imcrazytwkr/formdrain/repositories"
	"github.com/imcrazytwkr/formdrain/routes"
	"github.com/imcrazytwkr/formdrain/services"
)

type authRouter struct {
	sessionRepository repositories.SessionRepository
	accountService    services.AccountService
}

func NewAuthRouter(sessionRepository repositories.SessionRepository, accountService services.AccountService) routes.RouteContainer {
	return &authRouter{
		sessionRepository: sessionRepository,
		accountService:    accountService,
	}
}

func (r *authRouter) Router(router chi.Router) {
	router.Use(middleware.ContentTypeParser())
	router.Get("/login", r.HandleLoginForm)
	router.Post("/login", r.HandleLogin)
}
