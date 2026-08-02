package validators

import (
	"context"
	"net/netip"
)

type CaptchaValidator interface {
	Validate(ctx context.Context, responseToken string, hostname string, userIP netip.Addr) error
}
