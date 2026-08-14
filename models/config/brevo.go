package config

type BrevoConfig struct {
	SenderName  string `toml:"sender_name"`
	SenderEmail string `toml:"sender_email"`
}
