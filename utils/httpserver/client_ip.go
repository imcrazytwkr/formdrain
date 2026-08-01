package httpserver

import (
	"net"
	"net/http"
	"net/netip"
	"strings"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/utils/stringutil"
)

// RemoteIP returns the direct peer IP from RemoteAddr.
func RemoteIP(r *http.Request) netip.Addr {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	addr, err := netip.ParseAddr(host)
	if err != nil {
		return netip.Addr{}
	}

	return addr
}

// ClientIP returns the client IP, honoring X-Forwarded-For / X-Real-IP only when
// the direct peer is a trusted proxy (loopback, private, or link-local).
func ClientIP(r *http.Request) netip.Addr {
	remote := RemoteIP(r)
	if !remote.IsValid() || !isTrustedProxy(remote) {
		return remote
	}

	xff := r.Header.Get(constants.HeaderForwardedFor)
	if len(xff) > 0 {
		addr, err := netip.ParseAddr(strings.TrimSpace(stringutil.TakeUntilByte(xff, ',')))
		if err == nil {
			return addr
		}
	}

	xri := r.Header.Get(constants.HeaderRealIP)
	if len(xri) > 0 {
		addr, err := netip.ParseAddr(strings.TrimSpace(xri))
		if err == nil {
			return addr
		}
	}

	return remote
}

func isTrustedProxy(ip netip.Addr) bool {
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast()
}
