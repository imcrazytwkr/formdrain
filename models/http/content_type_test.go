package http_test

import (
	"testing"

	m "github.com/imcrazytwkr/formdrain/models/http"
)

func TestParseContentType(t *testing.T) {
	t.Parallel()

	if got := m.ParseContentType("application/json"); got != m.ContentTypeJSON {
		t.Fatalf("got %v", got)
	}
	if got := m.ParseContentType("text/html"); got != m.ContentTypeHTML {
		t.Fatalf("got %v", got)
	}
	if got := m.ParseContentType("text/plain"); got != m.ContentTypeUndefined {
		t.Fatalf("got %v", got)
	}
}

func TestParseFormContentType(t *testing.T) {
	t.Parallel()

	if got := m.ParseFormContentType("application/json"); got != m.ContentTypeJSON {
		t.Fatalf("got %v", got)
	}
	if got := m.ParseFormContentType("application/x-www-form-urlencoded"); got != m.ContentTypeHTML {
		t.Fatalf("got %v", got)
	}
	if got := m.ParseFormContentType("text/html"); got != m.ContentTypeUndefined {
		t.Fatalf("got %v", got)
	}
}

func TestContentType_String(t *testing.T) {
	t.Parallel()
	if m.ContentTypeJSON.String() != "application/json" {
		t.Fatal(m.ContentTypeJSON.String())
	}
}
