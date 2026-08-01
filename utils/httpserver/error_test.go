package httpserver_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/utils/httpserver"
	"github.com/imcrazytwkr/formdrain/validation"
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

func TestHandleValidationError_JSON(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	w.Header().Set(constants.HeaderContentType, m.ContentTypeJSON.String())

	joined := errors.Join(
		&validation.FieldError{Field: "email", Err: validation.ErrMissingRequiredField},
		&validation.FieldError{Field: "extra", Err: validation.ErrUnknownField},
	)
	httpserver.HandleValidationError(context.Background(), w, http.StatusBadRequest, joined)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["message"] != validation.ErrValidationFailed.Error() {
		t.Fatalf("message = %#v", body["message"])
	}
	errs, ok := body["errors"].(map[string]any)
	if !ok {
		t.Fatalf("errors = %#v", body["errors"])
	}
	if errs["email"] != validation.ErrMissingRequiredField.Error() {
		t.Fatalf("email = %#v", errs["email"])
	}
	if errs["extra"] != validation.ErrUnknownField.Error() {
		t.Fatalf("extra = %#v", errs["extra"])
	}
}
