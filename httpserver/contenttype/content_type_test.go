package contenttype_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/httpserver/contenttype"
	m "github.com/imcrazytwkr/formdrain/models/http"
)

func TestRequestContentType(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set(constants.HeaderContentType, "application/json; charset=utf-8")
	if got := contenttype.RequestContentType(r); got != "application/json" {
		t.Fatalf("RequestContentType = %q", got)
	}
}

func TestGetContentType(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	if got := contenttype.GetContentType(r); got != m.ContentTypeUndefined {
		t.Fatalf("GetContentType without ctx = %v", got)
	}

	ctx := contenttype.WithContentType(context.Background(), m.ContentTypeJSON)
	r = r.WithContext(ctx)
	if got := contenttype.GetContentType(r); got != m.ContentTypeJSON {
		t.Fatalf("GetContentType = %v", got)
	}
}
