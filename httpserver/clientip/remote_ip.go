package clientip

import (
	"net"
	"net/http"
	"net/netip"
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
