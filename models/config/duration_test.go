package config

import (
	"testing"
	"time"
)

func TestDuration_RoundTrip(t *testing.T) {
	t.Parallel()

	d := NewDuration(15 * time.Second)
	if d.Duration() != 15*time.Second {
		t.Fatalf("Duration() = %v", d.Duration())
	}

	text, err := d.MarshalText()
	if err != nil || string(text) != "15s" {
		t.Fatalf("MarshalText = %q err=%v", text, err)
	}

	var got Duration
	err = got.UnmarshalText([]byte("1h"))
	if err != nil || got.Duration() != time.Hour {
		t.Fatalf("UnmarshalText: got %v err=%v", got.Duration(), err)
	}
}

func TestDuration_UnmarshalInvalid(t *testing.T) {
	t.Parallel()

	var d Duration
	err := d.UnmarshalText([]byte("nope"))
	if err == nil {
		t.Fatal("expected error")
	}
}
