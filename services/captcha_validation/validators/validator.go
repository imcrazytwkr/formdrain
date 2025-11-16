package validators

import "context"

type CaptchaValidator interface {
	Validate(ctx context.Context, form map[string]any, hostname string, userIP string) error
}
