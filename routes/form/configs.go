package form

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imcrazytwkr/formdrain/utils/ginutil"
	"github.com/rs/zerolog"
)

func (r *formRouter) AttachFormConfigs(c *gin.Context) {
	ctx := c.Request.Context()
	log := zerolog.Ctx(ctx)

	formId := c.Param("formId")
	formConfig, err := r.formConfigRepository.GetFormConfigById(c.Request.Context(), formId)
	if err != nil {
		log.Err(err).Str("form_id", formId).Msg("error fetching form config")
		ginutil.HandleStatus(c, http.StatusInternalServerError)
		return
	}

	if formConfig == nil {
		ginutil.HandleError(c, http.StatusNotFound, getErrFormNotFound(formId))
		return
	}

	if len(formConfig.SiteId) < 1 {
		log.Error().Str("form_id", formId).Msg("form config has no sideId ref")
		ginutil.HandleStatus(c, http.StatusInternalServerError)
		return
	}

	siteConfig, err := r.siteConfigRepository.GetSiteConfigById(ctx, formConfig.SiteId.String())
	if err != nil {
		log.Err(err).Str("form_id", formId).Msg("error fetching site config for form")
		ginutil.HandleStatus(c, http.StatusInternalServerError)
		return
	}

	if siteConfig == nil {
		log.Error().Str("form_id", formId).Msg("could not find site config for form")
		ginutil.HandleStatus(c, http.StatusInternalServerError)
		return
	}

	c.Set(keyFormConfig, formConfig)
	c.Set(keySiteConfig, siteConfig)
	c.Next()
}
