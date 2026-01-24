package utils

func IsEmptyJSON(data []byte) bool {
	if len(data) == 0 {
		return true
	}

	switch string(data) {
	case "null", "[]", "{}", `""`:
		return true
	}

	return false
}
