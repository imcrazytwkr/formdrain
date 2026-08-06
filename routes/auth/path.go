package auth

import (
	"net/url"
	"strings"
)

// SanitizeAbsolutePath accepts only absolute paths (leading single '/').
// Absolute URLs, protocol-relative URLs, and relative paths become "/".
func SanitizeAbsolutePath(raw string) string {
	path := strings.TrimSpace(raw)
	if len(path) < 2 {
		// The only valid path of length 1 is root
		return "/"
	}

	if path[0] != '/' || path[1] == '/' {
		// Safe since we already know that length is >= 2
		return "/"
	}

	// Reject anything that url.Parse treats as having a scheme/host.
	u, err := url.Parse(path)
	if err != nil || len(u.Scheme) > 0 || len(u.Host) > 0 || len(u.Opaque) > 0 {
		return "/"
	}

	return path
}
