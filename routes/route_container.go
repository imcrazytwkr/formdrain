package routes

import "github.com/go-chi/chi/v5"

type RouteContainer interface {
	Router(router chi.Router)
}
