package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imcrazytwkr/formdrain/middleware"
	"github.com/rs/zerolog"
)

func TestLoggerWithConfig_LogsStatus(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := zerolog.New(&buf)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte("nope"))
	})

	h := middleware.LoggerWithConfig(&logger, nil)(next)
	req := httptest.NewRequest(http.MethodGet, "/x?y=1", nil)
	req.RemoteAddr = "203.0.113.1:9"
	h.ServeHTTP(httptest.NewRecorder(), req)

	var payload map[string]any
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if int(payload["status_code"].(float64)) != http.StatusTeapot {
		t.Fatalf("payload = %#v", payload)
	}
	if payload["path"] != "/x?y=1" {
		t.Fatalf("path = %#v", payload["path"])
	}
	if int(payload["body_size"].(float64)) != 4 {
		t.Fatalf("body_size = %#v", payload["body_size"])
	}
}

func TestLoggerWithConfig_SkipPath(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := zerolog.New(&buf)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := middleware.LoggerWithConfig(&logger, []string{"/health"})(next)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	h.ServeHTTP(httptest.NewRecorder(), req)

	if buf.Len() != 0 {
		t.Fatalf("expected no logs, got %s", buf.String())
	}
}
