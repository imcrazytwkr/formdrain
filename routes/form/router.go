package form

import (
	"github.com/gin-gonic/gin"
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

func (r *formRouter) Register(router gin.IRouter) {
	postRouter := router.Group("/:formId")
	postRouter.Use(middleware.ContentTypeParser())
	postRouter.Use(r.AttachFormConfigs)
	postRouter.Use(r.CheckCORS)
	postRouter.OPTIONS("/", r.HandlePreflight)
	postRouter.POST("/", r.HandleCreateForm)
}
