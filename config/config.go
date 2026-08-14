package config

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	m "github.com/imcrazytwkr/formdrain/models/config"
	"github.com/pelletier/go-toml/v2"
)

const (
	defaultPort              = 8080
	defaultShutdownTimeout   = 30 * time.Second
	defaultReadHeaderTimeout = 10 * time.Second
	defaultSessionTTL        = 24 * time.Hour
)

func Load(path string) (*m.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	return Parse(data)
}

func Parse(data []byte) (*m.Config, error) {
	cfg := defaults()
	dec := toml.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	err := dec.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	normalize(&cfg)

	err = validate(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func normalize(c *m.Config) {
	c.Server.Host = strings.TrimSpace(c.Server.Host)
	c.Auth.CookieDomain = strings.TrimSpace(c.Auth.CookieDomain)
	c.Notifiers.Discord.Username = strings.TrimSpace(c.Notifiers.Discord.Username)
	c.Notifiers.Discord.AvatarURL = strings.TrimSpace(c.Notifiers.Discord.AvatarURL)
	c.Notifiers.Brevo.SenderName = strings.TrimSpace(c.Notifiers.Brevo.SenderName)
	c.Notifiers.Brevo.SenderEmail = strings.TrimSpace(c.Notifiers.Brevo.SenderEmail)
}

func defaults() m.Config {
	return m.Config{
		Server: m.ServerConfig{
			Port:              defaultPort,
			ShutdownTimeout:   m.NewDuration(defaultShutdownTimeout),
			ReadHeaderTimeout: m.NewDuration(defaultReadHeaderTimeout),
		},
		Auth: m.AuthConfig{
			SessionTTL:     m.NewDuration(defaultSessionTTL),
			CookieSameSite: m.NewSameSite(http.SameSiteLaxMode),
		},
	}
}

func validate(c *m.Config) error {
	if len(c.Server.Host) > 0 && net.ParseIP(c.Server.Host) == nil {
		return fmt.Errorf("listen host %q is not a valid IP address", c.Server.Host)
	}

	if c.Server.Port < 0 || c.Server.Port > math.MaxInt16 {
		return fmt.Errorf("listen port number %d is invalid; Valid range is 0-%d", c.Server.Port, math.MaxInt16)
	}

	if c.Server.ShutdownTimeout.Duration() < time.Second {
		return fmt.Errorf("server.shutdown_timeout must be positive")
	}

	if c.Server.ReadHeaderTimeout.Duration() < time.Second {
		return fmt.Errorf("server.read_header_timeout must be positive")
	}

	if c.Auth.SessionTTL.Duration() < time.Second {
		return fmt.Errorf("auth.session_ttl must be positive")
	}

	if len(c.Notifiers.Discord.Username) < 1 {
		return fmt.Errorf("notifiers.discord.username is required")
	}

	if len(c.Notifiers.Brevo.SenderName) < 1 {
		return fmt.Errorf("notifiers.brevo.sender_name is required")
	}

	if len(c.Notifiers.Brevo.SenderEmail) < 1 {
		return fmt.Errorf("notifiers.brevo.sender_email is required")
	}

	return nil
}
