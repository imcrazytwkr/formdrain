package collate

func BoolToInt(value any) (int, bool) {
	v, ok := value.(bool)
	if !ok {
		return -1, false
	}

	if v {
		return 1, true
	}

	return 0, true
}
