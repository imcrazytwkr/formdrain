package httpserver

import (
	"net"
	"net/http"
	"strings"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/utils/stringutil"
)

// RemoteIP returns the direct peer IP from RemoteAddr.
func RemoteIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}

// ClientIP returns the client IP, honoring X-Forwarded-For / X-Real-IP only when
// the direct peer is a trusted proxy (loopback, private, or link-local).
func ClientIP(r *http.Request) string {
	remote := RemoteIP(r)
	if !isTrustedProxy(remote) {
		return remote
	}

	xff := r.Header.Get(constants.HeaderForwardedFor)
	if len(xff) > 0 {
		return strings.TrimSpace(stringutil.TakeUntilByte(xff, ','))
	}

	xri := r.Header.Get(constants.HeaderRealIP)
	if len(xri) > 0 {
		return strings.TrimSpace(xri)
	}

	return remote
}

func isTrustedProxy(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}

	return parsed.IsLoopback() || parsed.IsPrivate() || parsed.IsLinkLocalUnicast()
}
