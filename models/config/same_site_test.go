package config

import (
	"net/http"
	"testing"
)

func TestSameSite_RoundTrip(t *testing.T) {
	t.Parallel()

	cases := []struct {
		mode http.SameSite
		text string
	}{
		{http.SameSiteLaxMode, "lax"},
		{http.SameSiteStrictMode, "strict"},
		{http.SameSiteNoneMode, "none"},
	}

	for _, tc := range cases {
		s := NewSameSite(tc.mode)
		if s.SameSite() != tc.mode {
			t.Fatalf("%s: SameSite() = %v", tc.text, s.SameSite())
		}

		text, err := s.MarshalText()
		if err != nil || string(text) != tc.text {
			t.Fatalf("%s: MarshalText = %q err=%v", tc.text, text, err)
		}

		var got SameSite
		err = got.UnmarshalText([]byte(tc.text))
		if err != nil || got.SameSite() != tc.mode {
			t.Fatalf("%s: UnmarshalText got %v err=%v", tc.text, got.SameSite(), err)
		}
	}
}

func TestSameSite_UnmarshalInvalid(t *testing.T) {
	t.Parallel()

	var s SameSite
	err := s.UnmarshalText([]byte("weird"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSameSite_MarshalInvalid(t *testing.T) {
	t.Parallel()

	_, err := SameSite{}.MarshalText()
	if err == nil {
		t.Fatal("expected error")
	}
}
