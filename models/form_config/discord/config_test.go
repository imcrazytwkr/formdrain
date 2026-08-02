package discord_test

import (
	"testing"

	"github.com/imcrazytwkr/formdrain/models/common"
	"github.com/imcrazytwkr/formdrain/models/form_config/discord"
)

func TestRenderContent(t *testing.T) {
	t.Parallel()

	tmpl, err := common.NewTemplate("msg={{email}}")
	if err != nil {
		t.Fatal(err)
	}
	cfg := &discord.DiscordConfig{Template: tmpl}
	got, err := cfg.RenderContent(map[string]any{"email": "a@b.c"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "msg=a@b.c" {
		t.Fatalf("got %q", got)
	}
}

func TestRenderContent_NilTemplate(t *testing.T) {
	t.Parallel()

	cfg := &discord.DiscordConfig{}
	got, err := cfg.RenderContent(map[string]any{"email": "a@b.c"})
	if err != nil || got != "" {
		t.Fatalf("got %q err %v", got, err)
	}
}
