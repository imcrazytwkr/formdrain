package form

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/httpserver"
	"github.com/imcrazytwkr/formdrain/httpserver/clientip"
	"github.com/imcrazytwkr/formdrain/models/form_response"
	"github.com/imcrazytwkr/formdrain/utils/bodyparser"
	"github.com/imcrazytwkr/formdrain/utils/maputil"
	"github.com/imcrazytwkr/formdrain/validation"
)

func (r *formRouter) HandleCreateForm(w http.ResponseWriter, req *http.Request) {
	log := getLoggerForAction(req.Context(), actionSend)
	ctx := log.WithContext(req.Context())

	formData, err := bodyparser.Parse(req)
	if err != nil {
		if errors.Is(err, bodyparser.ErrBodyTooLarge) {
			httpserver.HandleError(ctx, w, http.StatusRequestEntityTooLarge, errFormTooLarge)
			return
		}

		httpserver.HandleError(ctx, w, http.StatusBadRequest, err)
		return
	}

	formId := chi.URLParam(req, "formId")

	formConfig, ok := getFormConfig(ctx)
	if !ok {
		log.Error().Msg("form config does not exist in context after being attached")
		httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
		return
	}

	siteConfig, ok := getSiteConfig(ctx)
	if !ok {
		log.Error().Msg("site config does not exist in context after being attached")
		httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
		return
	}

	clientIP := clientip.ClientIP(req)
	if !clientIP.IsValid() {
		log.Error().Msg("could not parse client IP")
		httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
		return
	}

	captchaField := formConfig.CaptchaTokenField()
	err = r.captchaValidationService.Validate(
		ctx,
		formConfig.CaptchaType,
		siteConfig.CaptchaSecret(formConfig.CaptchaType),
		maputil.String(formData, captchaField),
		siteConfig.Hostname,
		clientIP,
	)
	if err != nil {
		if errors.Is(err, constants.ErrCaptchaNotPassed) {
			httpserver.HandleError(ctx, w, http.StatusBadRequest, err)
			return
		}

		log.Err(err).
			Str("captcha_provider", formConfig.CaptchaType.String()).
			Str("form_id", formId).
			Str("user_ip", clientIP.String()).
			Msg("error validating user captcha")

		httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
		return
	}

	// Captcha tokens are not part of the field schema / stored payload.
	delete(formData, captchaField)
	payload, err := validation.ValidateFormPayload(formConfig.FieldSchema, formData)
	if err != nil {
		response := validation.NewValidationErrorResponse(http.StatusBadRequest, err)
		httpserver.HandleResponse(ctx, w, http.StatusBadRequest, "errors/validation", response)
		return
	}

	responseId, err := r.formResponseRepository.SaveFormResponse(ctx, &form_response.FormResponse{
		FormId:        formConfig.FormId,
		SchemaVersion: formConfig.SchemaVersion,
		ClientIP:      clientIP,
		Payload:       payload,
	})
	if err != nil {
		log.Err(err).Str("form_id", formId).Msg("could not save response into DB")
		httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
		return
	}

	log.Debug().Str("form_id", formId).Str("response_id", responseId).Msg("saved response")

	r.notificaionService.SendAsync(ctx, formConfig.Notifiers, payload)

	if len(formConfig.RedirectTo) > 0 {
		httpserver.HandleRedirect(ctx, w, http.StatusSeeOther, "form/redirect", formConfig.RedirectTo, nil)
		return
	}

	httpserver.HandleResponse(ctx, w, http.StatusOK, "form/success", map[string]any{
		"response_id": responseId,
	})
}
