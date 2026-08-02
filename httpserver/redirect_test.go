package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/httpserver"
	"github.com/imcrazytwkr/formdrain/models/common"
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

func TestHandleResponseTemplate_OwnerOverride(t *testing.T) {
	t.Parallel()

	owner, err := common.NewTemplate("<p>custom-{{v}}</p>")
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	httpserver.HandleResponseTemplate(
		t.Context(),
		w,
		http.StatusOK,
		owner,
		map[string]any{"v": "page"},
	)
	if !strings.Contains(w.Body.String(), "custom-page") {
		t.Fatalf("body = %q", w.Body.String())
	}
	if strings.Contains(w.Body.String(), "Form submitted") {
		t.Fatalf("expected owner template, got system body %q", w.Body.String())
	}
}

func TestHandleRedirectTemplate_OwnerOverride(t *testing.T) {
	t.Parallel()

	owner, err := common.NewTemplate(`<a href="{{redirect_to}}">continue</a>`)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	httpserver.HandleRedirectTemplate(
		t.Context(),
		w,
		http.StatusSeeOther,
		owner,
		"https://example.com/next",
		map[string]any{"redirect_to": "https://example.com/next"},
	)
	if loc := w.Header().Get(constants.HeaderLocation); loc != "https://example.com/next" {
		t.Fatalf("location = %q", loc)
	}
	body := w.Body.String()
	if !strings.Contains(body, "https://example.com/next") || !strings.Contains(body, "continue") {
		t.Fatalf("body = %q", body)
	}
	if strings.Contains(body, "Redirecting") {
		t.Fatalf("expected owner redirect template, got %q", body)
	}
}
