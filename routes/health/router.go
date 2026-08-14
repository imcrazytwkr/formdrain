package health

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/imcrazytwkr/formdrain/routes"
)

type healthRouter struct {
	db *sql.DB
}

func NewHealthRouter(db *sql.DB) routes.RouteContainer {
	return &healthRouter{db: db}
}

func (r *healthRouter) Router(router chi.Router) {
	router.Get("/livez", r.HandleLive)
	router.Get("/readyz", r.HandleReady)
}
