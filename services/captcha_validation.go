package services

import (
	"context"
	"net/netip"

	"github.com/imcrazytwkr/formdrain/models/common"
)

type CaptchaValidationService interface {
	Validate(
		ctx context.Context,
		captchaType common.CaptchaType,
		secret string,
		responseToken string,
		hostname string,
		userIP netip.Addr,
	) error
}
