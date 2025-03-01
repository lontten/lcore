package lcutils

import "strings"

func HasText(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) > 0
}

func HasTextP(s *string) bool {
	if s == nil {
		return false
	}
	var s2 = strings.TrimSpace(*s)
	return len(s2) > 0
}

func StrContainsAll(str string, list ...string) bool {
	for _, v := range list {
		if !strings.Contains(str, v) {
			return false
		}
	}
	return true
}

func StrContainsAny(str string, list ...string) bool {
	for _, v := range list {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
}
