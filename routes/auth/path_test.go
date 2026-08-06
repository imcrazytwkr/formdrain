package auth

import "testing"

func TestSanitizeAbsolutePath(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"", "/"},
		{"   ", "/"},
		{"/admin", "/admin"},
		{"/admin/settings", "/admin/settings"},
		{"https://evil.example/", "/"},
		{"http://evil.example/x", "/"},
		{"//evil", "/"},
		{"admin", "/"},
		{"./x", "/"},
		{"../x", "/"},
		{"/ok?x=1", "/ok?x=1"},
	}
	for _, tc := range cases {
		if got := SanitizeAbsolutePath(tc.in); got != tc.want {
			t.Errorf("SanitizeAbsolutePath(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
