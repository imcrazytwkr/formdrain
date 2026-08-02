package form_config_test

import (
	"encoding/json"
	"testing"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
)

func TestParseCaptchaType(t *testing.T) {
	t.Parallel()

	got, err := fc.ParseCaptchaType("hcaptcha")
	if err != nil || got != fc.CaptchaTypeHcaptcha {
		t.Fatalf("got %v err %v", got, err)
	}

	got, err = fc.ParseCaptchaType("recaptcha")
	if err != nil || got != fc.CaptchaTypeRecaptcha {
		t.Fatalf("got %v err %v", got, err)
	}

	_, err = fc.ParseCaptchaType("nope")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCaptchaType_JSONRoundTrip(t *testing.T) {
	t.Parallel()

	raw, err := json.Marshal(fc.CaptchaTypeHcaptcha)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != `"hcaptcha"` {
		t.Fatalf("marshal = %s", raw)
	}

	var got fc.CaptchaType
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got != fc.CaptchaTypeHcaptcha {
		t.Fatalf("got %v", got)
	}

	err = json.Unmarshal([]byte(`"bogus"`), &got)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCaptchaType_String(t *testing.T) {
	t.Parallel()
	if fc.CaptchaTypeRecaptcha.String() != "recaptcha" {
		t.Fatal(fc.CaptchaTypeRecaptcha.String())
	}
}
