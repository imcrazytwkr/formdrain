package config

import (
	"fmt"
	"net/http"
	"strings"
)

type SameSite struct {
	s http.SameSite
}

func NewSameSite(s http.SameSite) SameSite {
	return SameSite{s: s}
}

func (s SameSite) SameSite() http.SameSite {
	return s.s
}

func (s SameSite) MarshalText() ([]byte, error) {
	switch s.s {
	case http.SameSiteLaxMode:
		return []byte("lax"), nil
	case http.SameSiteStrictMode:
		return []byte("strict"), nil
	case http.SameSiteNoneMode:
		return []byte("none"), nil
	default:
		return nil, fmt.Errorf("same-site %d is invalid; expected lax, strict, or none", s.s)
	}
}

func (s *SameSite) UnmarshalText(text []byte) error {
	v := strings.ToLower(strings.TrimSpace(string(text)))
	switch v {
	case "lax":
		s.s = http.SameSiteLaxMode
		return nil
	case "strict":
		s.s = http.SameSiteStrictMode
		return nil
	case "none":
		s.s = http.SameSiteNoneMode
		return nil
	default:
		return fmt.Errorf("same-site %q is invalid; expected lax, strict, or none", v)
	}
}
