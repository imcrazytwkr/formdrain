package config_test

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/imcrazytwkr/formdrain/config"
	m "github.com/imcrazytwkr/formdrain/models/config"
)

const validTOML = `
[log]
mode = "release"

[server]
host = "127.0.0.1"
port = 3000
shutdown_timeout = "15s"
read_header_timeout = "5s"

[auth]
session_ttl = "1h"
cookie_secure = true
cookie_same_site = "strict"
cookie_domain = "example.com"

[notifiers.discord]
username = "FormDrain"
avatar_url = "https://example.com/avatar.png"

[notifiers.brevo]
sender_name = "Forms"
sender_email = "forms@example.com"
`

func TestParse_Valid(t *testing.T) {
	cfg, err := config.Parse([]byte(validTOML))
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Log.Mode != m.LogModeRelease {
		t.Fatalf("log.mode = %q", cfg.Log.Mode)
	}
	if cfg.Server.Host != "127.0.0.1" || cfg.Server.Port != 3000 {
		t.Fatalf("server = %+v", cfg.Server)
	}
	if cfg.Server.ShutdownTimeout.Duration() != 15*time.Second || cfg.Server.ReadHeaderTimeout.Duration() != 5*time.Second {
		t.Fatalf("timeouts = %v %v", cfg.Server.ShutdownTimeout.Duration(), cfg.Server.ReadHeaderTimeout.Duration())
	}
	if cfg.Auth.SessionTTL.Duration() != time.Hour || !cfg.Auth.CookieSecure || cfg.Auth.CookieSameSite.SameSite() != http.SameSiteStrictMode || cfg.Auth.CookieDomain != "example.com" {
		t.Fatalf("auth = %+v", cfg.Auth)
	}
	if cfg.Notifiers.Discord.Username != "FormDrain" || cfg.Notifiers.Discord.AvatarURL != "https://example.com/avatar.png" {
		t.Fatalf("discord = %+v", cfg.Notifiers.Discord)
	}
	if cfg.Notifiers.Brevo.SenderName != "Forms" || cfg.Notifiers.Brevo.SenderEmail != "forms@example.com" {
		t.Fatalf("brevo = %+v", cfg.Notifiers.Brevo)
	}
}

func TestParse_Defaults(t *testing.T) {
	cfg, err := config.Parse([]byte(`
[notifiers.discord]
username = "bot"

[notifiers.brevo]
sender_name = "Forms"
sender_email = "forms@example.com"
`))
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Log.Mode != m.LogModeTrace {
		t.Fatalf("log.mode = %q", cfg.Log.Mode)
	}
	if cfg.Server.Host != "" || cfg.Server.Port != 8080 {
		t.Fatalf("server = %+v", cfg.Server)
	}
	if cfg.Server.ShutdownTimeout.Duration() != 30*time.Second || cfg.Server.ReadHeaderTimeout.Duration() != 10*time.Second {
		t.Fatalf("timeouts = %v %v", cfg.Server.ShutdownTimeout.Duration(), cfg.Server.ReadHeaderTimeout.Duration())
	}
	if cfg.Auth.SessionTTL.Duration() != 24*time.Hour || cfg.Auth.CookieSecure || cfg.Auth.CookieSameSite.SameSite() != http.SameSiteLaxMode {
		t.Fatalf("auth = %+v", cfg.Auth)
	}
}

func TestParse_UnknownKey(t *testing.T) {
	_, err := config.Parse([]byte(`
extra = 1

[notifiers.discord]
username = "bot"

[notifiers.brevo]
sender_name = "Forms"
sender_email = "forms@example.com"
`))
	if err == nil {
		t.Fatal("expected unknown key error")
	}
}

func TestParse_InvalidHost(t *testing.T) {
	_, err := config.Parse([]byte(`
[server]
host = "not-an-ip"

[notifiers.discord]
username = "bot"

[notifiers.brevo]
sender_name = "Forms"
sender_email = "forms@example.com"
`))
	if err == nil || !strings.Contains(err.Error(), "valid IP") {
		t.Fatalf("err = %v", err)
	}
}

func TestParse_InvalidPort(t *testing.T) {
	_, err := config.Parse([]byte(`
[server]
port = 99999

[notifiers.discord]
username = "bot"

[notifiers.brevo]
sender_name = "Forms"
sender_email = "forms@example.com"
`))
	if err == nil || !strings.Contains(err.Error(), "port") {
		t.Fatalf("err = %v", err)
	}
}

func TestParse_InvalidDuration(t *testing.T) {
	_, err := config.Parse([]byte(`
[server]
shutdown_timeout = "nope"

[notifiers.discord]
username = "bot"

[notifiers.brevo]
sender_name = "Forms"
sender_email = "forms@example.com"
`))
	if err == nil || !strings.Contains(err.Error(), "invalid duration") {
		t.Fatalf("err = %v", err)
	}
}

func TestParse_InvalidSameSite(t *testing.T) {
	_, err := config.Parse([]byte(`
[auth]
cookie_same_site = "weird"

[notifiers.discord]
username = "bot"

[notifiers.brevo]
sender_name = "Forms"
sender_email = "forms@example.com"
`))
	if err == nil || !strings.Contains(err.Error(), "same-site") {
		t.Fatalf("err = %v", err)
	}
}

func TestParse_EmptyLogMode(t *testing.T) {
	cfg, err := config.Parse([]byte(`
[log]
mode = ""

[notifiers.discord]
username = "bot"

[notifiers.brevo]
sender_name = "Forms"
sender_email = "forms@example.com"
`))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Log.Mode != m.LogModeTrace {
		t.Fatalf("log.mode = %q", cfg.Log.Mode)
	}
}

func TestParse_InvalidLogMode(t *testing.T) {
	_, err := config.Parse([]byte(`
[log]
mode = "debug"

[notifiers.discord]
username = "bot"

[notifiers.brevo]
sender_name = "Forms"
sender_email = "forms@example.com"
`))
	if err == nil || !strings.Contains(err.Error(), "log mode") {
		t.Fatalf("err = %v", err)
	}
}

func TestParse_MissingIdentities(t *testing.T) {
	_, err := config.Parse([]byte(`
[notifiers.brevo]
sender_name = "Forms"
sender_email = "forms@example.com"
`))
	if err == nil || !strings.Contains(err.Error(), "username") {
		t.Fatalf("err = %v", err)
	}

	_, err = config.Parse([]byte(`
[notifiers.discord]
username = "bot"

[notifiers.brevo]
sender_email = "forms@example.com"
`))
	if err == nil || !strings.Contains(err.Error(), "sender_name") {
		t.Fatalf("err = %v", err)
	}

	_, err = config.Parse([]byte(`
[notifiers.discord]
username = "bot"

[notifiers.brevo]
sender_name = "Forms"
`))
	if err == nil || !strings.Contains(err.Error(), "sender_email") {
		t.Fatalf("err = %v", err)
	}
}

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	err := os.WriteFile(path, []byte(validTOML), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(path)
	if err != nil || cfg.Server.Port != 3000 {
		t.Fatalf("cfg=%#v err=%v", cfg, err)
	}

	_, err = config.Load(filepath.Join(dir, "missing.toml"))
	if err == nil {
		t.Fatal("expected missing file error")
	}
}

func TestParse_ExampleFile(t *testing.T) {
	data, err := os.ReadFile("../config.toml.example")
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := config.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Port != 8080 || cfg.Log.Mode != m.LogModeTrace || cfg.Notifiers.Discord.Username != "FormDrain" {
		t.Fatalf("example cfg = %+v", cfg)
	}
}

func TestParse_TrimsStrings(t *testing.T) {
	cfg, err := config.Parse([]byte(`
[server]
host = "  127.0.0.1  "

[auth]
cookie_domain = "  example.com  "

[notifiers.discord]
username = "  FormDrain  "
avatar_url = "  https://example.com/avatar.png  "

[notifiers.brevo]
sender_name = "  Forms  "
sender_email = "  forms@example.com  "
`))
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Server.Host != "127.0.0.1" {
		t.Fatalf("host = %q", cfg.Server.Host)
	}
	if cfg.Auth.CookieDomain != "example.com" {
		t.Fatalf("cookie_domain = %q", cfg.Auth.CookieDomain)
	}
	if cfg.Notifiers.Discord.Username != "FormDrain" || cfg.Notifiers.Discord.AvatarURL != "https://example.com/avatar.png" {
		t.Fatalf("discord = %+v", cfg.Notifiers.Discord)
	}
	if cfg.Notifiers.Brevo.SenderName != "Forms" || cfg.Notifiers.Brevo.SenderEmail != "forms@example.com" {
		t.Fatalf("brevo = %+v", cfg.Notifiers.Brevo)
	}
}

func TestParse_TrimmedBlankUsername(t *testing.T) {
	_, err := config.Parse([]byte(`
[notifiers.discord]
username = "   "

[notifiers.brevo]
sender_name = "Forms"
sender_email = "forms@example.com"
`))
	if err == nil || !strings.Contains(err.Error(), "username") {
		t.Fatalf("err = %v", err)
	}
}

func TestParse_ZeroDuration(t *testing.T) {
	_, err := config.Parse([]byte(`
[auth]
session_ttl = "0s"

[notifiers.discord]
username = "bot"

[notifiers.brevo]
sender_name = "Forms"
sender_email = "forms@example.com"
`))
	if err == nil || !strings.Contains(err.Error(), "session_ttl") {
		t.Fatalf("err = %v", err)
	}
}
