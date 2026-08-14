package config

import (
	"net"
	"strconv"
)

type ServerConfig struct {
	Host              string   `toml:"host"`
	Port              int      `toml:"port"`
	ShutdownTimeout   Duration `toml:"shutdown_timeout"`
	ReadHeaderTimeout Duration `toml:"read_header_timeout"`
}

func (c ServerConfig) Addr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}
