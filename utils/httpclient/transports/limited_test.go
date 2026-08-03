package transports_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/imcrazytwkr/formdrain/utils/httpclient/transports"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
	"golang.org/x/time/rate"
)

func TestLimitedTransport_CallsBase(t *testing.T) {
	t.Parallel()

	called := false
	base := testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
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

	synctest.Test(t, func(t *testing.T) {
		base := testutil.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			t.Fatal("base should not be called")
			return nil, nil
		})

		limiter := rate.NewLimiter(rate.Limit(1), 1)
		if err := limiter.Wait(t.Context()); err != nil {
			t.Fatal(err)
		}
		rt := transports.LimitedTransport(base, limiter)

		ctx, cancel := context.WithTimeout(t.Context(), time.Millisecond)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com/", nil)
		if err != nil {
			t.Fatal(err)
		}

		_, err = rt.RoundTrip(req)
		if err == nil {
			t.Fatal("expected context error")
		}
	})
}
