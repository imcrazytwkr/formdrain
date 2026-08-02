package httpserver_test

import (
	"testing"

	"github.com/imcrazytwkr/formdrain/httpserver"
)

func TestGetTemplate_UnknownPanics(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	_ = httpserver.GetTemplate("missing/nope")
}

func TestGetTemplate_Embedded(t *testing.T) {
	t.Parallel()

	tmpl := httpserver.GetTemplate("form/success")
	if tmpl == nil {
		t.Fatal("expected embedded template")
	}
}
