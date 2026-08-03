package discord_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/models/common"
	"github.com/imcrazytwkr/formdrain/models/form_config/discord"
	dn "github.com/imcrazytwkr/formdrain/services/notification/notifiers/discord"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func TestSend_Success(t *testing.T) {
	t.Parallel()

	var sawURL, sawCT string
	client := &http.Client{
		Transport: testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			sawURL = r.URL.String()
			sawCT = r.Header.Get(constants.HeaderContentType)
			body, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(body), "hello a@b.c") {
				t.Fatalf("body = %s", body)
			}
			return &http.Response{
				StatusCode: http.StatusNoContent,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	tmpl, err := common.NewTemplate("hello {{email}}")
	if err != nil {
		t.Fatal(err)
	}

	n := dn.NewDiscordNotifier("bot", "https://avatar", client)
	err = n.Send(t.Context(), &discord.DiscordConfig{
		Webhooks: []*discord.WebhookKey{{Snowflake: "123", Token: "tok"}},
		Title:    "t",
		Template: tmpl,
	}, map[string]any{"email": "a@b.c"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sawURL, "/webhooks/123/tok") {
		t.Fatalf("url = %q", sawURL)
	}
	if sawCT != constants.ContentTypeJson {
		t.Fatalf("Content-Type = %q", sawCT)
	}
}

func TestSend_SuccessStatusOK(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	n := dn.NewDiscordNotifier("bot", "avatar", client)
	err := n.Send(t.Context(), &discord.DiscordConfig{
		Webhooks: []*discord.WebhookKey{{Snowflake: "1", Token: "t"}},
		Template: mustTemplate(t, "hi"),
	}, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSend_NoWebhooks(t *testing.T) {
	t.Parallel()

	n := dn.NewDiscordNotifier("bot", "avatar", &http.Client{})
	err := n.Send(t.Context(), &discord.DiscordConfig{
		Webhooks: nil,
		Template: mustTemplate(t, "x"),
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSend_HTTPError(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadGateway,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	n := dn.NewDiscordNotifier("bot", "avatar", client)
	err := n.Send(t.Context(), &discord.DiscordConfig{
		Webhooks: []*discord.WebhookKey{{Snowflake: "1", Token: "t"}},
		Template: mustTemplate(t, "hi"),
	}, map[string]any{})
	if err == nil || !strings.Contains(err.Error(), "502") {
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
