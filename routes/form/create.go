package form

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/models/form_response"
	"github.com/imcrazytwkr/formdrain/utils/bodyparser"
	"github.com/imcrazytwkr/formdrain/utils/httpserver"
	"github.com/imcrazytwkr/formdrain/utils/logutil"
	"github.com/imcrazytwkr/formdrain/validation"
)

func (r *formRouter) HandleCreateForm(w http.ResponseWriter, req *http.Request) {
	log := getLoggerForAction(req.Context(), actionSend)
	ctx := log.WithContext(req.Context())

	contentType := httpserver.GetContentType(req)
	parser, ok := bodyparser.ParsersNew[contentType]
	if !ok {
		httpserver.HandleError(ctx, w, http.StatusUnsupportedMediaType, getErrUnsupportedFormType(contentType))
		return
	}

	if req.ContentLength > maxBodySize {
		httpserver.HandleError(ctx, w, http.StatusRequestEntityTooLarge, errFormTooLarge)
		return
	}

	// @NOTE: despite what you may think considering the earlier Content-Length check,
	// you can never trust your users to correctly set the headers
	body, err := io.ReadAll(io.LimitReader(req.Body, maxBodySize))
	if err != nil {
		httpserver.HandleError(ctx, w, http.StatusBadRequest, getErrMalformedFormData(contentType))
		return
	}

	formData, err := parser.Parse(body)
	if err != nil {
		httpserver.HandleError(ctx, w, http.StatusBadRequest, getErrMalformedFormData(contentType))
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

	clientIP := httpserver.ClientIP(req)
	if !clientIP.IsValid() {
		log.Error().Msg("could not parse client IP")
		httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
		return
	}

	err = r.captchaValidationService.Validate(ctx, formConfig.CaptchaType, formData, siteConfig.Hostname, clientIP)
	switch err {
	case nil:
		// Captcha check passed
		break
	case constants.ErrCaptchaNotPassed:
		httpserver.HandleError(ctx, w, http.StatusBadRequest, err)
		return
	default:
		log.Err(err).
			Str("captcha_provider", formConfig.CaptchaType.String()).
			Str("form_id", formId).
			Str("user_ip", clientIP.String()).
			Msg("error validating user captcha")

		httpserver.HandleStatus(ctx, w, http.StatusInternalServerError)
		return
	}

	// Captcha tokens are not part of the field schema / stored payload.
	delete(formData, "h-captcha")
	delete(formData, "g-recaptcha")

	payload, err := validation.ValidateFormPayload(formConfig.FieldSchema, formData)
	if err != nil {
		response := validation.NewValidationErrorResponse(http.StatusBadRequest, err)
		httpserver.HandleResponse(ctx, w, http.StatusBadRequest, "errors/validation.html", response)
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

	err = r.notificaionService.Send(formConfig.Notifiers, payload)
	if err != nil {
		logutil.UnwrapErr(log.Error(), err).Msg("failed to send notifications")
	}

	if len(formConfig.RedirectTo) > 0 {
		httpserver.HandleRedirect(ctx, w, http.StatusSeeOther, "form/redirect.html", formConfig.RedirectTo, nil)
		return
	}

	httpserver.HandleResponse(ctx, w, http.StatusOK, "form/success.html", map[string]any{
		"formConfig": formConfig,
		"siteConfig": siteConfig,
	})
}
