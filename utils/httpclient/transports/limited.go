package transports

import (
	"net/http"

	"golang.org/x/time/rate"
)

type limitedTransport struct {
	base    http.RoundTripper
	limiter *rate.Limiter
}

func LimitedTransport(base http.RoundTripper, limiter *rate.Limiter) http.RoundTripper {
	transport := &limitedTransport{
		base:    base,
		limiter: limiter,
	}

	if base == nil {
		transport.base = DefaultTransport()
	}

	return transport
}

func (t *limitedTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	err := t.limiter.Wait(request.Context())
	if err != nil {
		return nil, err
	}

	return t.base.RoundTrip(request)
}
