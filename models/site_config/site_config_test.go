package site_config

import (
	"testing"

	"github.com/imcrazytwkr/formdrain/models/common"
)

func TestCaptchaSecret(t *testing.T) {
	t.Parallel()

	cfg := &SiteConfig{
		HcaptchaSecret:  "h-secret",
		RecaptchaSecret: "r-secret",
	}

	if got := cfg.CaptchaSecret(common.CaptchaTypeHcaptcha); got != "h-secret" {
		t.Fatalf("hcaptcha = %q", got)
	}
	if got := cfg.CaptchaSecret(common.CaptchaTypeRecaptcha); got != "r-secret" {
		t.Fatalf("recaptcha = %q", got)
	}
	if got := cfg.CaptchaSecret(common.CaptchaTypeUndefined); got != "" {
		t.Fatalf("undefined = %q", got)
	}
	if got := (*SiteConfig)(nil).CaptchaSecret(common.CaptchaTypeHcaptcha); got != "" {
		t.Fatalf("nil = %q", got)
	}
}
