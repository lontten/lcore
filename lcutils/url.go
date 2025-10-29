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

func SanitizeFileName4URL(url string, size ...int) string {
	var sb strings.Builder
	sb.Grow(len(url)) // 预分配内存提高效率

	for _, c := range url {
		if _, ok := reservedChars[c]; ok {
			sb.WriteByte('_')
		} else {
			sb.WriteRune(c)
		}
	}

	s := sb.String()
	l := 15
	if len(size) > 0 {
		l = size[0]
	}

	runes := []rune(s)
	// 如果字符总数小于等于15，直接返回原字符串
	if len(runes) <= l {
		return s
	}
	// 否则返回前15个字符
	return string(runes[:l])
}
