package brevo_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/models/common"
	mcfg "github.com/imcrazytwkr/formdrain/models/config"
	"github.com/imcrazytwkr/formdrain/models/form_config/brevo"
	bn "github.com/imcrazytwkr/formdrain/services/notification/notifiers/brevo"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func TestSend_Success(t *testing.T) {
	t.Parallel()

	var sawURL string
	var sawAccept, sawCT, sawKey string
	client := &http.Client{
		Transport: testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			sawURL = r.URL.String()
			sawAccept = r.Header.Get(constants.HeaderAccept)
			sawCT = r.Header.Get(constants.HeaderContentType)
			sawKey = r.Header.Get("api-key")
			body, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(body), "body a@b.c") {
				t.Fatalf("body = %s", body)
			}
			return &http.Response{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	tmpl, err := common.NewTemplate("body {{email}}")
	if err != nil {
		t.Fatal(err)
	}

	n := bn.NewBrevoNotifier(mcfg.BrevoConfig{SenderName: "Sender", SenderEmail: "from@example.com"}, "test-api-key", client)
	err = n.Send(t.Context(), &brevo.BrevoConfig{
		Recipients: []*brevo.EmailContact{{Name: "A", Address: "a@b.c"}},
		Subject:    "subj",
		Template:   tmpl,
	}, map[string]any{"email": "a@b.c"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sawURL, "api.brevo.com") {
		t.Fatalf("url = %q", sawURL)
	}
	if sawAccept != constants.ContentTypeJson {
		t.Fatalf("Accept = %q", sawAccept)
	}
	if sawCT != constants.ContentTypeJson {
		t.Fatalf("Content-Type = %q", sawCT)
	}
	if sawKey != "test-api-key" {
		t.Fatalf("api-key = %q", sawKey)
	}
}

func TestSend_NoRecipients(t *testing.T) {
	t.Parallel()

	n := bn.NewBrevoNotifier(mcfg.BrevoConfig{SenderName: "Sender", SenderEmail: "from@example.com"}, "key", &http.Client{})
	err := n.Send(t.Context(), &brevo.BrevoConfig{
		Recipients: nil,
		Template:   mustTemplate(t, "x"),
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSend_NotConfigured(t *testing.T) {
	t.Parallel()

	n := bn.NewBrevoNotifier(mcfg.BrevoConfig{SenderName: "Sender", SenderEmail: "from@example.com"}, "", &http.Client{})
	err := n.Send(t.Context(), &brevo.BrevoConfig{
		Recipients: []*brevo.EmailContact{{Address: "a@b.c"}},
		Subject:    "s",
		Template:   mustTemplate(t, "hi"),
	}, map[string]any{})
	if !errors.Is(err, bn.ErrNotConfigured) {
		t.Fatalf("err = %v", err)
	}
}

func TestSend_HTTPError(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	n := bn.NewBrevoNotifier(mcfg.BrevoConfig{SenderName: "Sender", SenderEmail: "from@example.com"}, "key", client)
	err := n.Send(t.Context(), &brevo.BrevoConfig{
		Recipients: []*brevo.EmailContact{{Address: "a@b.c"}},
		Subject:    "s",
		Template:   mustTemplate(t, "hi"),
	}, map[string]any{})
	if err == nil || !strings.Contains(err.Error(), "500") {
		t.Fatalf("err = %v", err)
	}
}

func mustTemplate(t *testing.T, raw string) *common.Template {
	t.Helper()
	tmpl, err := common.NewTemplate(raw)
	if err != nil {
		t.Fatal(err)
	}
	return tmpl
}
