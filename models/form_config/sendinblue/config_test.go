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
	got := cfg.RenderContent(map[string]any{"email": "a@b.c"})
	if got != "body=a@b.c" {
		t.Fatalf("got %q", got)
	}
}
