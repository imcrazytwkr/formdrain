package origin_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/httpserver/origin"
)

func TestParseOriginHost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		host       string
		origin     string
		referer    string
		wantHost   string
		wantErr    error
		wantErrMsg string
	}{
		{
			name:     "no origin",
			host:     "api.example.com",
			wantHost: "",
		},
		{
			name:     "origin only",
			host:     "api.example.com",
			origin:   "https://forms.example.com",
			wantHost: "forms.example.com",
		},
		{
			name:     "matching referer",
			host:     "api.example.com",
			origin:   "https://forms.example.com",
			referer:  "https://forms.example.com/contact",
			wantHost: "forms.example.com",
		},
		{
			name:     "mismatched referer",
			host:     "api.example.com",
			origin:   "https://forms.example.com",
			referer:  "https://evil.example.com/x",
			wantHost: "",
			wantErr:  origin.ErrOriginMismatch,
		},
		{
			name:       "invalid referer",
			host:       "api.example.com",
			origin:     "https://forms.example.com",
			referer:    "://bad",
			wantHost:   "forms.example.com",
			wantErrMsg: "Referer header contains invalid URL",
		},
		{
			name:     "origin equals server host",
			host:     "api.example.com",
			origin:   "https://api.example.com",
			wantHost: "",
		},
		{
			name:       "invalid origin scheme",
			host:       "api.example.com",
			origin:     "ftp://forms.example.com",
			wantHost:   "",
			wantErrMsg: "Origin header contains invalid URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodPost, "https://api.example.com/form/1/", nil)
			r.Host = tt.host
			if tt.origin != "" {
				r.Header.Set(constants.HeaderOrigin, tt.origin)
			}
			if tt.referer != "" {
				r.Header.Set(constants.HeaderReferer, tt.referer)
			}

			got, err := origin.ParseOriginHost(r)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("err = %v, want %v", err, tt.wantErr)
				}
			} else if tt.wantErrMsg != "" {
				if err == nil || err.Error() != tt.wantErrMsg {
					t.Fatalf("err = %v, want message %q", err, tt.wantErrMsg)
				}
			} else if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}

			if got != tt.wantHost {
				t.Fatalf("host = %q, want %q", got, tt.wantHost)
			}
		})
	}
}
