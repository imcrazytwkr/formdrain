package origin

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/imcrazytwkr/formdrain/constants"
)

func errInvalidURL(header string) error {
	return fmt.Errorf("%s header contains invalid URL", header)
}

var ErrOriginMismatch = errors.New("origin and referer headers do not match")

func ParseOriginHost(r *http.Request) (string, error) {
	originHost, err := parseSourceUrl(r, constants.HeaderOrigin)
	if err != nil {
		return "", err
	}

	if len(originHost) < 1 {
		// Not a CORS request
		return "", nil
	}

	refererHost, err := parseSourceUrl(r, constants.HeaderReferer)
	if err != nil {
		return originHost, err
	}

	if len(refererHost) < 1 {
		// No Referer header
		return originHost, nil
	}

	if !strings.EqualFold(originHost, refererHost) {
		return "", ErrOriginMismatch
	}

	return originHost, nil
}

func parseSourceUrl(r *http.Request, header string) (string, error) {
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

	return sourceUrl.Host, nil
}
