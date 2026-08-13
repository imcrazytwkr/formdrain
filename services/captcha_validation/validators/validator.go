package validators

import (
	"context"
	"net/netip"
)

type CaptchaValidator interface {
	Validate(ctx context.Context, secret string, responseToken string, hostname string, userIP netip.Addr) error
}
