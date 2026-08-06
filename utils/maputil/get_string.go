package maputil

import "strings"

func GetString[K comparable](m map[K]any, key K) (string, bool) {
	v, ok := m[key]
	if !ok {
		return "", ok
	}

	s, ok := v.(string)
	return strings.TrimSpace(s), ok
}
