package httpserver_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/utils/httpserver"
)

func TestHandleResponse_String(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	httpserver.HandleResponse(context.Background(), w, http.StatusOK, "", "plain")
	if w.Body.String() != "plain" {
		t.Fatalf("body = %q", w.Body.String())
	}
}

func TestHandleResponse_JSON(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	w.Header().Set(constants.HeaderContentType, m.ContentTypeJSON.String())
	httpserver.HandleResponse(context.Background(), w, http.StatusOK, "", map[string]any{"ok": true})

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
	httpserver.HandleResponse(context.Background(), w, http.StatusOK, "form/success.html", nil)
	if !strings.Contains(w.Body.String(), "Form submitted") {
		t.Fatalf("body = %q", w.Body.String())
	}
}
