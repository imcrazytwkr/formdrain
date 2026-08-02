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
	got := cfg.RenderContent(map[string]any{"email": "a@b.c"})
	if got != "msg=a@b.c" {
		t.Fatalf("got %q", got)
	}
}
