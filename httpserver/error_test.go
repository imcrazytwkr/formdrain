package httpserver_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/httpserver"
	m "github.com/imcrazytwkr/formdrain/models/http"
)

func TestHandleError_JSON(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	w.Header().Set(constants.HeaderContentType, m.ContentTypeJSON.String())

	httpserver.HandleError(t.Context(), w, http.StatusBadRequest, errors.New("nope"))

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
	httpserver.HandleStatus(t.Context(), w, http.StatusNotFound)

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
	httpserver.HandleError(t.Context(), w, http.StatusInternalServerError, errors.New("secret"))

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["message"] != http.StatusText(http.StatusInternalServerError) {
		t.Fatalf("body = %#v", body)
	}
}

func TestHandleResponse_ValidationHTML(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	httpserver.HandleResponse(t.Context(), w, http.StatusBadRequest, "errors/validation", map[string]any{
		"status":  http.StatusBadRequest,
		"message": "form validation failed",
		"errors": []map[string]string{
			{"field": "email", "message": "required form field is missing"},
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
