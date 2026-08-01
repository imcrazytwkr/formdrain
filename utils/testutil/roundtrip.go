package testutil

import "net/http"

// RoundTripFunc adapts a function to http.RoundTripper for tests.
type RoundTripFunc func(*http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
