package captcha_validation

import (
	"context"
	"net/http"
	"net/netip"

	"github.com/imcrazytwkr/formdrain/models/common"
	"github.com/imcrazytwkr/formdrain/services"
	v "github.com/imcrazytwkr/formdrain/services/captcha_validation/validators"
	"github.com/imcrazytwkr/formdrain/services/captcha_validation/validators/hcaptcha"
	"github.com/imcrazytwkr/formdrain/services/captcha_validation/validators/recaptcha"
	"github.com/rs/zerolog"
)

type httpCaptchaValidationService struct {
	validators map[common.CaptchaType]v.CaptchaValidator
}

func NewHttpCaptchaValidationService(httpClient *http.Client, logger *zerolog.Logger) services.CaptchaValidationService {
	return &httpCaptchaValidationService{
		validators: map[common.CaptchaType]v.CaptchaValidator{
			common.CaptchaTypeHcaptcha:  hcaptcha.NewHcaptchaValidator(httpClient),
			common.CaptchaTypeRecaptcha: recaptcha.NewRecaptchaValidator(httpClient),
		},
	}
}

func (s *httpCaptchaValidationService) Validate(
	ctx context.Context,
	captchaType common.CaptchaType,
	secret string,
	responseToken string,
	hostname string,
	userIP netip.Addr,
) error {
	validator, ok := s.validators[captchaType]
	if !ok {
		return getErrNoCaptchaImpl(captchaType)
	}

	if len(secret) < 1 {
		return getErrMissingCaptchaSecret(captchaType)
	}

	return validator.Validate(ctx, secret, responseToken, hostname, userIP)
}
