package recaptcha_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/services/captcha_validation/validators/recaptcha"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestValidate_MissingToken(t *testing.T) {
	t.Parallel()

	v := recaptcha.NewRecaptchaValidator("secret", &http.Client{})
	err := v.Validate(context.Background(), "", "example.com", netip.Addr{})
	if !errors.Is(err, recaptcha.ErrNoRecaptchaToken) {
		t.Fatalf("err = %v", err)
	}
}

func TestValidate_Success(t *testing.T) {
	t.Parallel()

	var sawRequest *http.Request
	client := &http.Client{
		Transport: testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			sawRequest = r
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			values, err := url.ParseQuery(string(body))
			if err != nil {
				t.Fatal(err)
			}
			if values.Get("secret") != "sec" || values.Get("response") != "tok" || values.Get("remoteip") != "1.2.3.4" {
				t.Fatalf("payload = %v", values)
			}
			if r.Header.Get(constants.HeaderContentType) != constants.ContentTypeForm {
				t.Fatalf("content-type = %q", r.Header.Get(constants.HeaderContentType))
			}
			return jsonResponse(http.StatusOK, `{"success":true,"hostname":"example.com"}`), nil
		}),
	}

	v := recaptcha.NewRecaptchaValidator("sec", client)
	err := v.Validate(context.Background(), "tok", "example.com", netip.MustParseAddr("1.2.3.4"))
	if err != nil {
		t.Fatal(err)
	}
	if sawRequest == nil || sawRequest.Method != http.MethodPost {
		t.Fatalf("request = %#v", sawRequest)
	}
}

func TestValidate_SuccessWithBOM(t *testing.T) {
	t.Parallel()

	body := string([]byte{0xEF, 0xBB, 0xBF}) + `{"success":true,"hostname":"example.com"}`
	client := &http.Client{
		Transport: testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			return jsonResponse(http.StatusOK, body), nil
		}),
	}

	v := recaptcha.NewRecaptchaValidator("sec", client)
	if err := v.Validate(context.Background(), "tok", "example.com", netip.Addr{}); err != nil {
		t.Fatal(err)
	}
}

func TestValidate_NotPassed(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			return jsonResponse(http.StatusOK, `{"success":false}`), nil
		}),
	}

	v := recaptcha.NewRecaptchaValidator("sec", client)
	err := v.Validate(context.Background(), "tok", "example.com", netip.Addr{})
	if !errors.Is(err, constants.ErrCaptchaNotPassed) {
		t.Fatalf("err = %v", err)
	}
}

func TestValidate_HostnameMismatch(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			return jsonResponse(http.StatusOK, `{"success":true,"hostname":"evil.com"}`), nil
		}),
	}

	v := recaptcha.NewRecaptchaValidator("sec", client)
	err := v.Validate(context.Background(), "tok", "example.com", netip.Addr{})
	if !errors.Is(err, constants.ErrCaptchaNotPassed) {
		t.Fatalf("err = %v", err)
	}
}

func TestValidate_EmptyHostnameAllowed(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			return jsonResponse(http.StatusOK, `{"success":true}`), nil
		}),
	}

	v := recaptcha.NewRecaptchaValidator("sec", client)
	if err := v.Validate(context.Background(), "tok", "example.com", netip.Addr{}); err != nil {
		t.Fatal(err)
	}
}

func TestValidate_NonOKStatus(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			return jsonResponse(http.StatusBadGateway, `oops`), nil
		}),
	}

	v := recaptcha.NewRecaptchaValidator("sec", client)
	err := v.Validate(context.Background(), "tok", "example.com", netip.Addr{})
	if err == nil || !strings.Contains(err.Error(), "502") {
		t.Fatalf("err = %v", err)
	}
}

func TestValidate_MalformedJSON(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			return jsonResponse(http.StatusOK, `{`), nil
		}),
	}

	v := recaptcha.NewRecaptchaValidator("sec", client)
	err := v.Validate(context.Background(), "tok", "example.com", netip.Addr{})
	if err == nil || !strings.Contains(err.Error(), "malformed body") {
		t.Fatalf("err = %v", err)
	}
}
