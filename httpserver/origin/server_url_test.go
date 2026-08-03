package origin_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/httpserver/origin"
)

func TestParseServerURL_Forwarded(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	r.Host = "localhost:8080"
	r.RemoteAddr = "127.0.0.1:9999"
	r.Header.Set(constants.HeaderForwardedFor, "198.51.100.1")
	r.Header.Set(constants.HeaderForwardedHost, "public.example.com")
	r.Header.Set(constants.HeaderForwardedPort, "443")

	u, err := origin.ParseServerURL(r)
	if err != nil {
		t.Fatal(err)
	}
	if u.Host != "public.example.com" || u.Scheme != "https" {
		t.Fatalf("got %#v", u)
	}

	host, err := origin.ParseServerHost(r)
	if err != nil || host != "public.example.com" {
		t.Fatalf("ParseServerHost = %q, %v", host, err)
	}

	r.Header.Set(constants.HeaderForwardedPort, "nope")
	_, err = origin.ParseServerURL(r)
	if !errors.Is(err, origin.ErrInvalidForwardedPort) {
		t.Fatalf("err = %v, want ErrInvalidForwardedPort", err)
	}
}
