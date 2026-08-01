package transports_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/imcrazytwkr/formdrain/utils/httpclient/transports"
	"golang.org/x/time/rate"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestLimitedTransport_CallsBase(t *testing.T) {
	t.Parallel()

	called := false
	base := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		called = true
		return &http.Response{
			StatusCode: http.StatusNoContent,
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    r,
		}, nil
	})

	limiter := rate.NewLimiter(rate.Inf, 1)
	rt := transports.LimitedTransport(base, limiter)

	req, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if !called {
		t.Fatal("base RoundTrip was not called")
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("status = %d", resp.StatusCode)
	}
}

func TestLimitedTransport_ContextCancel(t *testing.T) {
	t.Parallel()

	base := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		t.Fatal("base should not be called")
		return nil, nil
	})

	limiter := rate.NewLimiter(rate.Limit(1), 1)
	_ = limiter.Wait(context.Background()) // consume burst
	rt := transports.LimitedTransport(base, limiter)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = rt.RoundTrip(req)
	if err == nil {
		t.Fatal("expected context error")
	}
}
