package httpserver

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/imcrazytwkr/formdrain/constants"
)

var ErrInvalidForwardedPort = errors.New(constants.HeaderForwardedFor + " header is invalid")

func ParseServerHost(r *http.Request) (string, error) {
	serverUrl, err := ParseServerURL(r)
	if serverUrl != nil {
		return serverUrl.Host, err
	}

	return "", err
}

func ParseServerURL(r *http.Request) (*url.URL, error) {
	result := &url.URL{}

	result.Host = r.Host
	if len(result.Host) < 1 {
		result.Host = r.URL.Host
	}

	result.Scheme = r.URL.Scheme

	// @NOTE: if this condition holds, we are running behind a trusted proxy
	if ClientIP(r) != RemoteIP(r) {
		forwardedUrl, err := parseForwardedHost(r)
		if err != nil {
			// Returning error here because these headers take prevalence
			return nil, err
		}

		if forwardedUrl != nil {
			if len(forwardedUrl.Host) > 0 {
				result.Host = forwardedUrl.Host
			}

			if len(forwardedUrl.Scheme) > 0 {
				result.Scheme = forwardedUrl.Scheme
			}
		}
	}

	if len(result.Host) < 1 {
		return nil, nil
	}

	return result, nil
}

func parseForwardedHost(r *http.Request) (*url.URL, error) {
	result := &url.URL{}

	headerHost := r.Header.Get(constants.HeaderForwardedHost)
	if len(headerHost) < 1 {
		return nil, nil
	}

	rawPort := r.Header.Get(constants.HeaderForwardedPort)
	if len(rawPort) < 1 {
		result.Host = headerHost
		return result, nil
	}

	port, err := strconv.Atoi(rawPort)
	if err != nil {
		return nil, ErrInvalidForwardedPort
	}

	switch port {
	case 80:
		result.Host = headerHost
		result.Scheme = "http"
	case 443:
		result.Host = headerHost
		result.Scheme = "https"
	default:
		result.Host = fmt.Sprintf("%s:%d", headerHost, port)
	}

	return result, nil
}
