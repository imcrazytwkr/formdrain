package config

import (
	"strings"
	"time"
)

type Duration struct {
	d time.Duration
}

func NewDuration(d time.Duration) Duration {
	return Duration{d: d}
}

func (d Duration) Duration() time.Duration {
	return d.d
}

func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.d.String()), nil
}

func (d *Duration) UnmarshalText(text []byte) error {
	parsed, err := time.ParseDuration(strings.TrimSpace(string(text)))
	if err != nil {
		return err
	}

	d.d = parsed
	return nil
}
