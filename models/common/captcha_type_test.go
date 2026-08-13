package common_test

import (
	"encoding/json"
	"testing"

	"github.com/imcrazytwkr/formdrain/models/common"
)

func TestParseCaptchaType(t *testing.T) {
	t.Parallel()

	got, err := common.ParseCaptchaType("hcaptcha")
	if err != nil || got != common.CaptchaTypeHcaptcha {
		t.Fatalf("got %v err %v", got, err)
	}

	got, err = common.ParseCaptchaType("recaptcha")
	if err != nil || got != common.CaptchaTypeRecaptcha {
		t.Fatalf("got %v err %v", got, err)
	}

	_, err = common.ParseCaptchaType("nope")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCaptchaType_JSONRoundTrip(t *testing.T) {
	t.Parallel()

	raw, err := json.Marshal(common.CaptchaTypeHcaptcha)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != `"hcaptcha"` {
		t.Fatalf("marshal = %s", raw)
	}

	var got common.CaptchaType
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got != common.CaptchaTypeHcaptcha {
		t.Fatalf("got %v", got)
	}

	err = json.Unmarshal([]byte(`"bogus"`), &got)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCaptchaType_String(t *testing.T) {
	t.Parallel()
	if common.CaptchaTypeRecaptcha.String() != "recaptcha" {
		t.Fatal(common.CaptchaTypeRecaptcha.String())
	}
}

func TestCaptchaType_DefaultTokenField(t *testing.T) {
	t.Parallel()
	if got := common.CaptchaTypeHcaptcha.DefaultTokenField(); got != "h-captcha" {
		t.Fatalf("hcaptcha default = %q", got)
	}
	if got := common.CaptchaTypeRecaptcha.DefaultTokenField(); got != "g-recaptcha" {
		t.Fatalf("recaptcha default = %q", got)
	}
	if got := common.CaptchaTypeUndefined.DefaultTokenField(); got != "" {
		t.Fatalf("undefined default = %q", got)
	}
}
