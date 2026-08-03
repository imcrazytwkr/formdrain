package contenttype_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/httpserver/contenttype"
	m "github.com/imcrazytwkr/formdrain/models/http"
)

func TestNegotiate(t *testing.T) {
	t.Parallel()

	offers := []m.ContentType{m.ContentTypeHTML, m.ContentTypeJSON}

	tests := []struct {
		name   string
		accept string
		want   m.ContentType
	}{
		{name: "empty", accept: "", want: m.ContentTypeUndefined},
		{name: "json", accept: "application/json", want: m.ContentTypeJSON},
		{name: "html", accept: "text/html", want: m.ContentTypeHTML},
		{name: "q values", accept: "application/json;q=0.5, text/html;q=0.9", want: m.ContentTypeHTML},
		{name: "no match", accept: "image/png", want: m.ContentTypeUndefined},
		{name: "star", accept: "*/*", want: m.ContentTypeHTML},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.accept != "" {
				r.Header.Set(constants.HeaderAccept, tt.accept)
			}
			if got := contenttype.Negotiate(r, offers); got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
