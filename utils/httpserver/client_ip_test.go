package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/utils/httpserver"
)

func TestClientIP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		remoteAddr string
		xff        string
		xri        string
		want       string
		wantRemote string
	}{
		{
			name:       "direct public peer ignores forwarded",
			remoteAddr: "203.0.113.10:1234",
			xff:        "198.51.100.1",
			want:       "203.0.113.10",
			wantRemote: "203.0.113.10",
		},
		{
			name:       "trusted loopback uses xff",
			remoteAddr: "127.0.0.1:1234",
			xff:        "198.51.100.1, 10.0.0.1",
			want:       "198.51.100.1",
			wantRemote: "127.0.0.1",
		},
		{
			name:       "trusted private uses x-real-ip",
			remoteAddr: "10.0.0.2:1234",
			xri:        "198.51.100.9",
			want:       "198.51.100.9",
			wantRemote: "10.0.0.2",
		},
		{
			name:       "trusted without headers uses remote",
			remoteAddr: "127.0.0.1:1234",
			want:       "127.0.0.1",
			wantRemote: "127.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.RemoteAddr = tt.remoteAddr
			if tt.xff != "" {
				r.Header.Set(constants.HeaderForwardedFor, tt.xff)
			}
			if tt.xri != "" {
				r.Header.Set(constants.HeaderRealIP, tt.xri)
			}

			if got := httpserver.ClientIP(r); got != tt.want {
				t.Fatalf("ClientIP = %q, want %q", got, tt.want)
			}
			if got := httpserver.RemoteIP(r); got != tt.wantRemote {
				t.Fatalf("RemoteIP = %q, want %q", got, tt.wantRemote)
			}
		})
	}
}
