package captcha_validation

import (
	"io"
	"net/http"
	"net/netip"
	"strings"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/models/common"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
	"github.com/rs/zerolog"
)

func TestValidate_UnknownType(t *testing.T) {
	t.Parallel()

	svc := NewHttpCaptchaValidationService(&http.Client{}, &zerolog.Logger{})
	err := svc.Validate(t.Context(), common.CaptchaTypeUndefined, "secret", "", "example.com", netip.Addr{})
	if err == nil || !strings.Contains(err.Error(), "catcha type") {
		t.Fatalf("err = %v", err)
	}
}

func TestValidate_HcaptchaHappyPath(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"success":true,"hostname":"example.com"}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	svc := NewHttpCaptchaValidationService(client, &zerolog.Logger{})
	err := svc.Validate(
		t.Context(),
		common.CaptchaTypeHcaptcha,
		"secret",
		"tok",
		"example.com",
		netip.MustParseAddr("1.2.3.4"),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidate_HcaptchaNotPassed(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"success":false}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	svc := NewHttpCaptchaValidationService(client, &zerolog.Logger{})
	err := svc.Validate(
		t.Context(),
		common.CaptchaTypeHcaptcha,
		"secret",
		"tok",
		"example.com",
		netip.Addr{},
	)
	if err != constants.ErrCaptchaNotPassed {
		t.Fatalf("err = %v", err)
	}
}

func TestValidate_MissingSecret(t *testing.T) {
	t.Parallel()

	svc := NewHttpCaptchaValidationService(&http.Client{}, &zerolog.Logger{})
	err := svc.Validate(t.Context(), common.CaptchaTypeHcaptcha, "", "tok", "example.com", netip.Addr{})
	if err == nil || !strings.Contains(err.Error(), "missing captcha secret") {
		t.Fatalf("err = %v", err)
	}
}
