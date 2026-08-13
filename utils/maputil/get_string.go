package maputil

import "strings"

func GetString[K comparable](m map[K]any, key K) (string, bool) {
	v, ok := m[key]
	if !ok {
		return "", false
	}

	s, ok := v.(string)
	if !ok {
		return "", false
	}

	return strings.TrimSpace(s), true
}

func String[K comparable](m map[K]any, key K) string {
	v, _ := GetString(m, key)
	return v
}
