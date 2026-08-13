package form_config_test

import (
	"testing"

	"github.com/imcrazytwkr/formdrain/models/common"
	fc "github.com/imcrazytwkr/formdrain/models/form_config"
)

func TestFormConfig_CaptchaTokenField(t *testing.T) {
	t.Parallel()

	cfg := &fc.FormConfig{CaptchaType: common.CaptchaTypeHcaptcha}
	if got := cfg.CaptchaTokenField(); got != "h-captcha" {
		t.Fatalf("default = %q", got)
	}

	cfg.CaptchaField = "cf-turnstile-response"
	if got := cfg.CaptchaTokenField(); got != "cf-turnstile-response" {
		t.Fatalf("custom = %q", got)
	}

	if got := (*fc.FormConfig)(nil).CaptchaTokenField(); got != "" {
		t.Fatalf("nil config = %q", got)
	}
}
