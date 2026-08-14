package config

type AuthConfig struct {
	SessionTTL     Duration `toml:"session_ttl"`
	CookieSecure   bool     `toml:"cookie_secure"`
	CookieSameSite SameSite `toml:"cookie_same_site"`
	CookieDomain   string   `toml:"cookie_domain"`
}
