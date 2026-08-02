package httpserver_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/httpserver"
	m "github.com/imcrazytwkr/formdrain/models/http"
)

func TestHandleResponse_String(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	httpserver.HandleResponse(t.Context(), w, http.StatusOK, "", "plain")
	if w.Body.String() != "plain" {
		t.Fatalf("body = %q", w.Body.String())
	}
}

func TestHandleResponse_JSON(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	w.Header().Set(constants.HeaderContentType, m.ContentTypeJSON.String())
	httpserver.HandleResponse(t.Context(), w, http.StatusOK, "", map[string]any{"ok": true})

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["ok"] != true {
		t.Fatalf("body = %#v", body)
	}
}

func TestHandleResponse_HTML(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	httpserver.HandleResponse(t.Context(), w, http.StatusOK, "form/success", nil)
	if !strings.Contains(w.Body.String(), "Form submitted") {
		t.Fatalf("body = %q", w.Body.String())
	}
}

func TestHandleResponse_JSONMarshalError(t *testing.T) {
	w := httptest.NewRecorder()
	w.Header().Set(constants.HeaderContentType, m.ContentTypeJSON.String())
	httpserver.HandleResponse(t.Context(), w, http.StatusOK, "", make(chan int))
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if w.Body.Len() != 0 {
		t.Fatalf("expected empty body on marshal failure, got %q", w.Body.String())
	}
}

func TestHandleResponse_MissingHTMLTemplate(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for unknown embedded template")
		}
	}()

	w := httptest.NewRecorder()
	httpserver.HandleResponse(t.Context(), w, http.StatusOK, "missing/nope", nil)
}

func TestGetTemplate_OwnerNilFallsBack(t *testing.T) {
	t.Parallel()

	tmpl := httpserver.GetTemplate("form/success")
	if tmpl == nil {
		t.Fatal("expected embedded template")
	}

	w := httptest.NewRecorder()
	httpserver.HandleResponseTemplate(t.Context(), w, http.StatusOK, tmpl, map[string]any{})
	if !strings.Contains(w.Body.String(), "Form submitted") {
		t.Fatalf("body = %q", w.Body.String())
	}
}

func TestHandleError_NonStandardStatus(t *testing.T) {
	w := httptest.NewRecorder()
	w.Header().Set(constants.HeaderContentType, m.ContentTypeJSON.String())
	httpserver.HandleStatus(t.Context(), w, 599)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", w.Code)
	}
}
