package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/middleware"
	"github.com/rs/zerolog"
)

func TestRequestId_Generates(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zerolog.Ctx(r.Context()).Info().Msg("inside")
		w.WriteHeader(http.StatusNoContent)
	})

	h := middleware.RequestId()(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "203.0.113.1:1234"
	req = req.WithContext(logger.WithContext(req.Context()))
	h.ServeHTTP(httptest.NewRecorder(), req)

	var payload map[string]any
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	id, _ := payload["request_id"].(string)
	if len(id) < 1 {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestRequestId_HonorsForwarded(t *testing.T) {
	t.Parallel()

	const want = "req-from-proxy"
	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zerolog.Ctx(r.Context()).Info().Msg("inside")
		w.WriteHeader(http.StatusNoContent)
	})

	h := middleware.RequestId()(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	req.Header.Set(constants.HeaderForwardedFor, "198.51.100.1")
	req.Header.Set(constants.HeaderRequestID, want)
	req = req.WithContext(logger.WithContext(req.Context()))
	h.ServeHTTP(httptest.NewRecorder(), req)

	var payload map[string]any
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["request_id"] != want {
		t.Fatalf("payload = %#v", payload)
	}
}
