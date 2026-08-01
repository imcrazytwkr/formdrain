package httpserver_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/utils/httpserver"
)

func TestRequestContentType(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set(constants.HeaderContentType, "application/json; charset=utf-8")
	if got := httpserver.RequestContentType(r); got != "application/json" {
		t.Fatalf("RequestContentType = %q", got)
	}
}

func TestGetContentType(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	if got := httpserver.GetContentType(r); got != m.ContentTypeUndefined {
		t.Fatalf("GetContentType without ctx = %v", got)
	}

	ctx := httpserver.WithContentType(context.Background(), m.ContentTypeJSON)
	r = r.WithContext(ctx)
	if got := httpserver.GetContentType(r); got != m.ContentTypeJSON {
		t.Fatalf("GetContentType = %v", got)
	}
}
