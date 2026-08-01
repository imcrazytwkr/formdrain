package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/middleware"
	m "github.com/imcrazytwkr/formdrain/models/http"
)

func TestResponseFormatParser_SetsHeader(t *testing.T) {
	t.Parallel()

	var got string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = w.Header().Get(constants.HeaderContentType)
	})

	h := middleware.ResponseFormatParser(m.ContentTypeHTML, m.ContentTypeJSON)(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	h.ServeHTTP(httptest.NewRecorder(), req)

	if got != m.ContentTypeJSON.String() {
		t.Fatalf("content-type = %q", got)
	}
}

func TestResponseFormatParser_LeavesUnset(t *testing.T) {
	t.Parallel()

	var got string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = w.Header().Get(constants.HeaderContentType)
	})

	h := middleware.ResponseFormatParser(m.ContentTypeHTML, m.ContentTypeJSON)(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(httptest.NewRecorder(), req)

	if got != "" {
		t.Fatalf("content-type = %q, want empty", got)
	}
}
