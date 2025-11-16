package maputil

func GetString[K comparable](obj map[K]any, key K) (string, bool) {
	raw, ok := obj[key]
	if !ok {
		return "", ok
	}

	str, ok := raw.(string)
	return str, ok
}
