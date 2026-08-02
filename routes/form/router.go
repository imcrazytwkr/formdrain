package form

import (
	"github.com/go-chi/chi/v5"
	"github.com/imcrazytwkr/formdrain/middleware"
	"github.com/imcrazytwkr/formdrain/repositories"
	"github.com/imcrazytwkr/formdrain/routes"
	"github.com/imcrazytwkr/formdrain/services"
)

type formRouter struct {
	formConfigRepository     repositories.FormConfigRepository
	siteConfigRepository     repositories.SiteConfigRepository
	formResponseRepository   repositories.FormResponseRepository
	captchaValidationService services.CaptchaValidationService
	notificaionService       services.NotificationService
}

func NewFormRouter(
	formConfigRepository repositories.FormConfigRepository,
	siteConfigRepository repositories.SiteConfigRepository,
	formResponseRepository repositories.FormResponseRepository,
	captchaValidationService services.CaptchaValidationService,
	notificaionService services.NotificationService,
) routes.RouteContainer {
	return &formRouter{
		formConfigRepository:     formConfigRepository,
		siteConfigRepository:     siteConfigRepository,
		formResponseRepository:   formResponseRepository,
		captchaValidationService: captchaValidationService,
		notificaionService:       notificaionService,
	}
}

func (r *formRouter) Router(router chi.Router) {
	router.Route("/{formId}", func(route chi.Router) {
		route.Use(middleware.ContentTypeParser())
		route.Use(r.AttachFormConfigs)
		route.Use(r.CheckCORS)
		route.Options("/", r.HandlePreflight)
		route.Post("/", r.HandleCreateForm)
	})
}
