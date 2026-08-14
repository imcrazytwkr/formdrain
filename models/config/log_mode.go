package config

import (
	"fmt"
	"strings"
)

type LogMode int8

const (
	LogModeTrace LogMode = iota
	LogModeRelease
)

func (m LogMode) String() string {
	switch m {
	case LogModeTrace:
		return "trace"
	case LogModeRelease:
		return "release"
	default:
		return ""
	}
}

func (m LogMode) MarshalText() ([]byte, error) {
	switch m {
	case LogModeTrace, LogModeRelease:
		return []byte(m.String()), nil
	default:
		return nil, fmt.Errorf("log mode %d is invalid; expected trace or release", m)
	}
}

func (m *LogMode) UnmarshalText(text []byte) error {
	v := strings.TrimSpace(string(text))
	switch v {
	case "", "trace":
		*m = LogModeTrace
		return nil
	case "release":
		*m = LogModeRelease
		return nil
	default:
		return fmt.Errorf("log mode %q is invalid; expected trace or release", v)
	}
}
