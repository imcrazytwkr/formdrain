package form

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/imcrazytwkr/formdrain/utils/httpserver"
	"github.com/rs/zerolog"
)

func (r *formRouter) AttachFormConfigs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		log := zerolog.Ctx(ctx)

		formId := chi.URLParam(req, "formId")
		formConfig, err := r.formConfigRepository.GetFormConfigById(ctx, formId)
		if err != nil {
			log.Err(err).Str("form_id", formId).Msg("error fetching form config")
			httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
			return
		}

		if formConfig == nil {
			httpserver.HandleError(ctx, w, http.StatusNotFound, getErrFormNotFound(formId))
			return
		}

		if formConfig.SiteId < 1 {
			log.Error().Str("form_id", formId).Msg("form config has no siteId ref")
			httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
			return
		}

		siteConfig, err := r.siteConfigRepository.GetSiteConfigById(ctx, strconv.FormatInt(formConfig.SiteId, 10))
		if err != nil {
			log.Err(err).Str("form_id", formId).Msg("error fetching site config for form")
			httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
			return
		}

		if siteConfig == nil {
			log.Error().Str("form_id", formId).Msg("could not find site config for form")
			httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
			return
		}

		ctx = withFormConfig(ctx, formConfig)
		ctx = withSiteConfig(ctx, siteConfig)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
