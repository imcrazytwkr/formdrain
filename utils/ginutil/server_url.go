package ginutil

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imcrazytwkr/formdrain/constants"
)

var errInvalidForwardedPort = errors.New("X-Forwarded-Port header is invalid")

func ParseServerHost(c *gin.Context) (string, error) {
	serverUrl, err := ParseServerURL(c)
	if serverUrl != nil {
		return serverUrl.Host, err
	}

	return "", err
}

func ParseServerURL(c *gin.Context) (*url.URL, error) {
	result := &url.URL{}

	result.Host = c.Request.Host
	if len(result.Host) < 1 {
		result.Host = c.Request.URL.Host
	}

	result.Scheme = c.Request.URL.Scheme

	// @NOTE: if this condition holds, we are running behind a trusted proxy
	if c.ClientIP() != c.RemoteIP() {
		forwardedUrl, err := parseForwardedHost(c)
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

func parseForwardedHost(c *gin.Context) (*url.URL, error) {
	result := &url.URL{}

	headerHost := c.GetHeader(constants.HeaderForwardedHost)
	if len(headerHost) < 1 {
		return nil, nil
	}

	rawPort := c.GetHeader(constants.HeaderForwardedPort)
	if len(rawPort) < 1 {
		result.Host = headerHost
		return result, nil
	}

	port, err := strconv.Atoi(rawPort)
	if err != nil {
		return nil, errInvalidForwardedPort
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
