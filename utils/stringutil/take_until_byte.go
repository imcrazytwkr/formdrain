package stringutil

import "strings"

func TakeUntilByte(s string, c byte) string {
	i := strings.IndexByte(s, c)
	if i < 0 {
		return s
	}

	if i == 0 {
		return ""
	}

	return s[:i]
}
