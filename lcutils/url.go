package lcutils

import (
	"strings"
)

// 定义RFC 3986规定的保留字符集合
var reservedChars = map[rune]struct{}{
	':': {}, '/': {}, '?': {}, '#': {}, '[': {}, ']': {},
	'@': {}, '!': {}, '$': {}, '&': {}, '\'': {}, '(': {},
	')': {}, '*': {}, '+': {}, ',': {}, ';': {}, '=': {},
}

func SanitizeFileName4URL(url string) string {
	var sb strings.Builder
	sb.Grow(len(url)) // 预分配内存提高效率

	for _, c := range url {
		if _, ok := reservedChars[c]; ok {
			sb.WriteByte('_')
		} else {
			sb.WriteRune(c)
		}
	}

	return sb.String()
}
