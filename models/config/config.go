package config

type Config struct {
	Log       LogConfig       `toml:"log"`
	Server    ServerConfig    `toml:"server"`
	Auth      AuthConfig      `toml:"auth"`
	Notifiers NotifiersConfig `toml:"notifiers"`
}
