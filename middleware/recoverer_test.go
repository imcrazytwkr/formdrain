package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/middleware"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/rs/zerolog"
)

func TestRecoverer_JSONBody(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	h := middleware.Recoverer()(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	rec.Header().Set(constants.HeaderContentType, m.ContentTypeJSON.String())
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if int(body["status"].(float64)) != http.StatusInternalServerError {
		t.Fatalf("body = %#v", body)
	}
	if body["message"] != http.StatusText(http.StatusInternalServerError) {
		t.Fatalf("body = %#v", body)
	}
}

func TestRecoverer_HTMLFallback(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	h := middleware.Recoverer()(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", rec.Code)
	}
	if ct := rec.Header().Get(constants.HeaderContentType); ct != m.ContentTypeHTML.String() {
		t.Fatalf("content-type = %q", ct)
	}
	if !strings.Contains(rec.Body.String(), "Internal Server Error") {
		t.Fatalf("body = %q", rec.Body.String())
	}
}

func TestRecoverer_LogsPanic(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("logged-boom")
	})

	h := middleware.Recoverer()(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(logger.WithContext(req.Context()))
	h.ServeHTTP(httptest.NewRecorder(), req)

	out := buf.String()
	if !strings.Contains(out, "panic recovered") {
		t.Fatalf("log = %s", out)
	}
	if !strings.Contains(out, "logged-boom") {
		t.Fatalf("log missing panic value: %s", out)
	}
	if !strings.Contains(out, "stack") {
		t.Fatalf("log missing stack: %s", out)
	}
}

func TestRecoverer_RePanicsErrAbortHandler(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(http.ErrAbortHandler)
	})

	h := middleware.Recoverer()(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	repanicked := false
	func() {
		defer func() {
			rvr := recover()
			if rvr == http.ErrAbortHandler {
				repanicked = true
			}
		}()
		h.ServeHTTP(httptest.NewRecorder(), req)
	}()

	if !repanicked {
		t.Fatal("expected http.ErrAbortHandler to be re-panicked")
	}
}

func TestRecoverer_UpgradeSkipsResponse(t *testing.T) {
	t.Parallel()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("upgrade-boom")
	})

	h := middleware.Recoverer()(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Connection", "Upgrade")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Body.Len() != 0 {
		t.Fatalf("expected empty body, got %q", rec.Body.String())
	}
	if rec.Code != http.StatusOK {
		// httptest defaults to 200 when WriteHeader was never called
		t.Fatalf("status = %d", rec.Code)
	}
}
