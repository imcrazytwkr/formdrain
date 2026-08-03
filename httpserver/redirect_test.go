package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/httpserver"
	m "github.com/imcrazytwkr/formdrain/models/http"
)

func TestHandleRedirect_JSONLocationOnly(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	w.Header().Set(constants.HeaderContentType, m.ContentTypeJSON.String())
	httpserver.HandleRedirect(t.Context(), w, http.StatusSeeOther, "form/redirect", "https://example.com/ok", nil)

	if loc := w.Header().Get(constants.HeaderLocation); loc != "https://example.com/ok" {
		t.Fatalf("location = %q", loc)
	}
	if w.Body.Len() != 0 {
		t.Fatalf("expected empty body, got %q", w.Body.String())
	}
}

func TestHandleRedirect_HTMLBody(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	httpserver.HandleRedirect(t.Context(), w, http.StatusSeeOther, "form/redirect", "https://example.com/ok", nil)

	if loc := w.Header().Get(constants.HeaderLocation); loc != "https://example.com/ok" {
		t.Fatalf("location = %q", loc)
	}
	if !strings.Contains(w.Body.String(), "Redirecting") {
		t.Fatalf("body = %q", w.Body.String())
	}
}
