package captcha_validation

import (
	"context"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/services"
	v "github.com/imcrazytwkr/formdrain/services/captcha_validation/validators"
	"github.com/imcrazytwkr/formdrain/services/captcha_validation/validators/hcaptcha"
	"github.com/imcrazytwkr/formdrain/services/captcha_validation/validators/recaptcha"
	"github.com/imcrazytwkr/formdrain/types"
	"github.com/rs/zerolog"
)

type httpCaptchaValidationService struct {
	validators map[fc.CaptchaType]v.CaptchaValidator
}

func NewHttpCaptchaValidationService(httpClient types.HttpClient, logger *zerolog.Logger) services.CaptchaValidationService {
	return &httpCaptchaValidationService{
		validators: map[fc.CaptchaType]v.CaptchaValidator{
			fc.CaptchaTypeHcaptcha:  hcaptcha.NewHcaptchaValidator("hcaptcha_secret", httpClient),
			fc.CaptchaTypeRecaptcha: recaptcha.NewRecaptchaValidator("recaptcha_secret", httpClient),
		},
	}
}

func (s *httpCaptchaValidationService) Validate(
	ctx context.Context,
	captchaType fc.CaptchaType,
	form map[string]any,
	hostname string,
	userIP string,
) error {
	validator, ok := s.validators[captchaType]
	if !ok {
		return getErrNoCaptchaImpl(captchaType)
	}

	return validator.Validate(ctx, form, hostname, userIP)
}
