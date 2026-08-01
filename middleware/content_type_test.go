package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/middleware"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/utils/httpserver"
)

func TestContentTypeParser_StoresAndGuesses(t *testing.T) {
	t.Parallel()

	var (
		stored m.ContentType
		respCT string
	)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stored = httpserver.GetContentType(r)
		respCT = w.Header().Get(constants.HeaderContentType)
		w.WriteHeader(http.StatusNoContent)
	})

	h := middleware.ContentTypeParser()(next)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(constants.HeaderContentType, "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if stored != m.ContentTypeJSON {
		t.Fatalf("stored = %v", stored)
	}
	if respCT != m.ContentTypeJSON.String() {
		t.Fatalf("response content-type = %q", respCT)
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestContentTypeParser_RejectsMissing(t *testing.T) {
	t.Parallel()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	h := middleware.ContentTypeParser()(next)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if called {
		t.Fatal("next should not be called")
	}
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestContentTypeParser_DoesNotOverrideAccept(t *testing.T) {
	t.Parallel()

	var respCT string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respCT = w.Header().Get(constants.HeaderContentType)
	})

	h := middleware.ResponseFormatParser(m.ContentTypeHTML, m.ContentTypeJSON)(
		middleware.ContentTypeParser()(next),
	)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(constants.HeaderAccept, "text/html")
	req.Header.Set(constants.HeaderContentType, "application/json")
	h.ServeHTTP(httptest.NewRecorder(), req)

	if respCT != m.ContentTypeHTML.String() {
		t.Fatalf("content-type = %q", respCT)
	}
}

func TestContentTypeParser_SkipsGet(t *testing.T) {
	t.Parallel()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})

	h := middleware.ContentTypeParser()(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !called || rec.Code != http.StatusNoContent {
		t.Fatalf("called=%v status=%d", called, rec.Code)
	}
}
