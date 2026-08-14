package config

type NotifiersConfig struct {
	Discord DiscordConfig `toml:"discord"`
	Brevo   BrevoConfig   `toml:"brevo"`
}
