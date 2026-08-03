package form

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/imcrazytwkr/formdrain/httpserver"
	"github.com/rs/zerolog"
)

func (r *formRouter) AttachFormConfigs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		log := zerolog.Ctx(ctx)

		formIdRaw := chi.URLParam(req, "formId")
		formId, err := strconv.ParseInt(formIdRaw, 10, 64)
		if err != nil || formId < 1 {
			httpserver.HandleError(ctx, w, http.StatusNotFound, getErrFormNotFound(formIdRaw))
			return
		}

		formConfig, err := r.formConfigRepository.GetFormConfigById(ctx, formId)
		if err != nil {
			log.Err(err).Int64("form_id", formId).Msg("error fetching form config")
			httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
			return
		}

		if formConfig == nil {
			httpserver.HandleError(ctx, w, http.StatusNotFound, getErrFormNotFound(formIdRaw))
			return
		}

		if formConfig.SiteId < 1 {
			log.Error().Int64("form_id", formId).Msg("form config has no siteId ref")
			httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
			return
		}

		siteConfig, err := r.siteConfigRepository.GetSiteConfigById(ctx, formConfig.SiteId)
		if err != nil {
			log.Err(err).Int64("form_id", formId).Msg("error fetching site config for form")
			httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
			return
		}

		if siteConfig == nil {
			log.Error().Int64("form_id", formId).Msg("could not find site config for form")
			httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
			return
		}

		ctx = withFormConfig(ctx, formConfig)
		ctx = withSiteConfig(ctx, siteConfig)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
