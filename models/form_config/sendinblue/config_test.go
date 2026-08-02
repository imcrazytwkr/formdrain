package sendinblue_test

import (
	"testing"

	"github.com/imcrazytwkr/formdrain/models/common"
	"github.com/imcrazytwkr/formdrain/models/form_config/sendinblue"
)

func TestRenderContent(t *testing.T) {
	t.Parallel()

	tmpl, err := common.NewTemplate("body={{email}}")
	if err != nil {
		t.Fatal(err)
	}
	cfg := &sendinblue.SendinblueConfig{Template: tmpl}
	got, err := cfg.RenderContent(map[string]any{"email": "a@b.c"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "body=a@b.c" {
		t.Fatalf("got %q", got)
	}
}

func TestRenderContent_NilTemplate(t *testing.T) {
	t.Parallel()

	cfg := &sendinblue.SendinblueConfig{}
	got, err := cfg.RenderContent(map[string]any{"email": "a@b.c"})
	if err != nil || got != "" {
		t.Fatalf("got %q err %v", got, err)
	}
}
