package validators

import (
	"context"
	"net/netip"
)

type CaptchaValidator interface {
	Validate(ctx context.Context, form map[string]any, hostname string, userIP netip.Addr) error
}
