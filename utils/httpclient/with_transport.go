package httpclient

import "net/http"

func WithTransport(client *http.Client, transport http.RoundTripper) *http.Client {
	return &http.Client{
		// Do no evil
		Transport: transport,
		Timeout:   max(client.Timeout, defaultRequestTimeout),

		// Safe to pass-through as-is
		CheckRedirect: client.CheckRedirect,
		Jar:           client.Jar,
	}
}
