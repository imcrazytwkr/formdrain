package services

import (
	"context"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
)

type CaptchaValidationService interface {
	Validate(ctx context.Context, captchaType fc.CaptchaType, form map[string]any, hostname string, userIP string) error
}
