package transports

import (
	"net"
	"net/http"
	"time"
)

const defaultConnectionTimeout = 30 * time.Second
const defaultConnectionKeepAlive = defaultConnectionTimeout

func defaultDialer() *net.Dialer {
	return &net.Dialer{
		Timeout:   defaultConnectionKeepAlive,
		KeepAlive: defaultConnectionTimeout,
	}
}

const defaultMaxIdleConns = 100
const defaultIdleTimeout = 90 * time.Second
const defaultHandshakeTimeout = 10 * time.Second
const defaultExpectContinueTimeout = time.Second

// Instantiates default transport so that it's safe to mutate as needed
func DefaultTransport() http.RoundTripper {
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           defaultDialer().DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          defaultMaxIdleConns,
		IdleConnTimeout:       defaultIdleTimeout,
		TLSHandshakeTimeout:   defaultHandshakeTimeout,
		ExpectContinueTimeout: defaultExpectContinueTimeout,
	}
}
