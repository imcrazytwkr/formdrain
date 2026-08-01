package httpserver

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/imcrazytwkr/formdrain/constants"
)

func errInvalidURL(header string) error {
	return fmt.Errorf("%s header contains invalid URL", header)
}

var ErrOriginMismatch = errors.New("origin and referer headers do not match")

func ParseOriginHost(r *http.Request) (string, error) {
	serverHost, _ := ParseServerHost(r)

	originHost, err := parseSourceUrl(r, constants.HeaderOrigin, serverHost)
	if err != nil {
		return "", err
	}

	if len(originHost) < 1 {
		// Not a CORS request
		return "", nil
	}

	refererHost, err := parseSourceUrl(r, constants.HeaderReferer, serverHost)
	if err != nil {
		return originHost, err
	}

	if len(refererHost) < 1 {
		// No Referer header
		return originHost, nil
	}

	if originHost != refererHost {
		return "", ErrOriginMismatch
	}

	return originHost, nil
}

func parseSourceUrl(r *http.Request, header string, serverHost string) (string, error) {
	source := r.Header.Get(header)
	if len(source) < 1 {
		return "", nil
	}

	sourceUrl, err := url.Parse(source)
	if err != nil {
		return "", errInvalidURL(header)
	}

	if len(sourceUrl.Scheme) > 0 && sourceUrl.Scheme != "http" && sourceUrl.Scheme != "https" {
		return "", errInvalidURL(header)
	}

	if len(serverHost) > 0 && sourceUrl.Host == serverHost {
		// Direct request from misconfigured Fetch
		return "", nil
	}

	return sourceUrl.Host, nil
}
