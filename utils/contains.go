package utils

func ArrayContainsString(h []string, s string) bool {
	for _, v := range h {
		if v == s {
			return true
		}
	}
	return false
}
