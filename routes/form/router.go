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

func (r *formRouter) Register(router chi.Router) {
	router.Route("/{formId}", func(rtr chi.Router) {
		rtr.Use(middleware.ContentTypeParser())
		rtr.Use(r.AttachFormConfigs)
		rtr.Use(r.CheckCORS)
		rtr.Options("/", r.HandlePreflight)
		rtr.Post("/", r.HandleCreateForm)
	})
}
