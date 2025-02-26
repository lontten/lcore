package utils

func StrNilToStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
