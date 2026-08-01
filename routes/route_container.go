package routes

import "github.com/go-chi/chi/v5"

type RouteContainer interface {
	Register(router chi.Router)
}
