package form

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/utils/bodyparser"
	"github.com/imcrazytwkr/formdrain/utils/ginutil"
	"github.com/imcrazytwkr/formdrain/utils/logutil"
)

func (r *formRouter) HandleCreateForm(c *gin.Context) {
	log := getLoggerForAction(c.Request.Context(), actionSend)
	ctx := log.WithContext(c.Request.Context())

	contentType := ginutil.GetContentType(c)
	parser, ok := bodyparser.ParsersNew[contentType]
	if !ok {
		ginutil.HandleError(c, http.StatusUnsupportedMediaType, getErrUnsupportedFormType(contentType))
		return
	}

	if c.Request.ContentLength > maxBodySize {
		ginutil.HandleError(c, http.StatusRequestEntityTooLarge, errFormTooLarge)
		return
	}

	// @NOTE: despite what you may think considering the earlier Content-Length check,
	// you can never trust your users to correctly set the headers
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, maxBodySize))
	if err != nil {
		ginutil.HandleError(c, http.StatusBadRequest, getErrMalformedFormData(contentType))
		return
	}

	form, err := parser.Parse(body)
	if err != nil {
		ginutil.HandleError(c, http.StatusBadRequest, getErrMalformedFormData(contentType))
		return
	}

	formId := c.Param("formId")

	formConfig, ok := getFormConfig(c)
	if !ok {
		log.Error().Msg("form config does not exist in context after being attached")
		ginutil.HandleStatus(c, http.StatusInternalServerError)
		return
	}

	siteConfig, ok := getSiteConfig(c)
	if !ok {
		log.Error().Msg("site config does not exist in context after being attached")
		ginutil.HandleStatus(c, http.StatusInternalServerError)
		return
	}

	userIP := c.ClientIP()
	if len(userIP) < 1 {
		log.Error().Msg("could not parse client IP")
		ginutil.HandleStatus(c, http.StatusInternalServerError)
		return
	}

	err = r.captchaValidationService.Validate(ctx, formConfig.CaptchaType, form, siteConfig.Hostname, userIP)
	switch err {
	case nil:
		// Captcha check passed
		break
	case constants.ErrCaptchaNotPassed:
		ginutil.HandleError(c, http.StatusBadRequest, err)
		return
	default:
		log.Err(err).
			Str("captcha_provider", formConfig.CaptchaType.String()).
			Str("form_id", formId).
			Str("user_ip", userIP).
			Msg("error validating user captcha")

		ginutil.HandleStatus(c, http.StatusInternalServerError)
		return
	}

	responseId, err := r.formResponseRepository.SaveFormResponse(ctx, form)
	if err != nil {
		log.Err(err).Str("form_id", formId).Msg("could not save response into DB")
		return
	}

	log.Debug().Str("form_id", formId).Str("response_id", responseId).Msg("saved response")

	err = r.notificaionService.Send(formConfig.Notifiers, form)
	if err != nil {
		logutil.UnwrapErr(log.Error(), err).Msg("failed to send notifications")
	}

	if len(formConfig.RedirectTo) > 0 {
		ginutil.HandleRedirect(c, http.StatusSeeOther, "form/redirect.html", formConfig.RedirectTo, nil)
		return
	}

	ginutil.HandleResponse(c, http.StatusOK, "form/success.html", &gin.H{
		"formConfig": formConfig,
		"siteConfig": siteConfig,
	})
}
