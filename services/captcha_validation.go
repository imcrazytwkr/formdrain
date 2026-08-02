package services

import (
	"context"
	"net/netip"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
)

type CaptchaValidationService interface {
	Validate(ctx context.Context, captchaType fc.CaptchaType, responseToken string, hostname string, userIP netip.Addr) error
}
