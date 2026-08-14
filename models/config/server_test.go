package config

import (
	"testing"
)

func TestServerConfig_Addr(t *testing.T) {
	t.Parallel()

	cfg := ServerConfig{Host: "127.0.0.1", Port: 3000}
	if got := cfg.Addr(); got != "127.0.0.1:3000" {
		t.Fatalf("Addr() = %q", got)
	}

	cfg = ServerConfig{Port: 8080}
	if got := cfg.Addr(); got != ":8080" {
		t.Fatalf("empty host Addr() = %q", got)
	}
}
