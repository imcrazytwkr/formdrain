package httpserver_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/utils/httpserver"
)

func TestHandleError_JSON(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	w.Header().Set(constants.HeaderContentType, m.ContentTypeJSON.String())

	httpserver.HandleError(context.Background(), w, http.StatusBadRequest, errors.New("nope"))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["message"] != "nope" || int(body["status"].(float64)) != http.StatusBadRequest {
		t.Fatalf("body = %#v", body)
	}
}

func TestHandleError_HTMLFallbackWhenUnset(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	httpserver.HandleStatus(context.Background(), w, http.StatusNotFound)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d", w.Code)
	}
	if ct := w.Header().Get(constants.HeaderContentType); ct != m.ContentTypeHTML.String() {
		t.Fatalf("content-type = %q", ct)
	}
	if !strings.Contains(w.Body.String(), "Not Found") {
		t.Fatalf("body = %q", w.Body.String())
	}
}

func TestHandleError_ServerHidesMessage(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	w.Header().Set(constants.HeaderContentType, m.ContentTypeJSON.String())
	httpserver.HandleError(context.Background(), w, http.StatusInternalServerError, errors.New("secret"))

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["message"] != http.StatusText(http.StatusInternalServerError) {
		t.Fatalf("body = %#v", body)
	}
}

func TestHandleResponse_ValidationHTML(t *testing.T) {
	err := httpserver.LoadTemplatesFromPath(filepath.Join("..", "..", "templates"))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	httpserver.HandleResponse(context.Background(), w, http.StatusBadRequest, "errors/validation.html", map[string]any{
		"Status":  http.StatusBadRequest,
		"Message": "form validation failed",
		"Errors": map[string]string{
			"email": "required form field is missing",
		},
	})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "form validation failed") {
		t.Fatalf("missing message: %q", body)
	}
	if !strings.Contains(body, "email") || !strings.Contains(body, "required form field is missing") {
		t.Fatalf("missing field errors: %q", body)
	}
}
