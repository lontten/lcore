package lcutils

import (
	"strings"
	"unicode"
)

// 定义RFC 3986规定的保留字符集合
var reservedChars = map[rune]struct{}{
	':': {}, '/': {}, '?': {}, '#': {}, '[': {}, ']': {},
	'@': {}, '!': {}, '$': {}, '&': {}, '\'': {}, '(': {},
	')': {}, '*': {}, '+': {}, ',': {}, ';': {}, '=': {},
}

func SanitizeFileName4URL(url string, size ...int) string {
	url = strings.TrimSpace(url)
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
	s = CleanString(s)
	s = processUnderscoresOptimized(s)
	runes := []rune(s)
	l := 0
	if len(size) > 0 {
		l = size[0]
	}

	// 如果字符总数小于等于15，直接返回原字符串
	if l == 0 || len(runes) <= l {
		return s
	}
	// 否则返回前15个字符
	return string(runes[:l])
}

// CleanString 将字符串中的空格、空字符和不可见字符替换为下划线
func CleanString(s string) string {
	return strings.Map(func(r rune) rune {
		// 检查字符是否为空格、空字符或不可见字符
		if unicode.IsSpace(r) || isInvisibleControlCharacter(r) {
			return '_' // 替换为下划线
		}
		return r // 保留原字符
	}, s)
}

// isInvisibleControlCharacter 检查字符是否为不可见的控制字符
func isInvisibleControlCharacter(r rune) bool {
	// 控制字符的 Unicode 范围是 0x0000 到 0x001F 和 0x007F 到 0x009F
	return (r >= 0x0000 && r <= 0x001F) || (r >= 0x007F && r <= 0x009F)
}

func processUnderscoresOptimized(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	var builder strings.Builder

	// 找到第一个非下划线字符的位置
	start := 0
	for start < len(runes) && runes[start] == '_' {
		start++
	}

	// 找到最后一个非下划线字符的位置
	end := len(runes) - 1
	for end >= 0 && runes[end] == '_' {
		end--
	}

	if start > end {
		return ""
	}

	// 处理中间部分
	for i := start; i <= end; i++ {
		current := runes[i]

		if current == '_' {
			// 只有当上一个字符不是下划线时才添加
			if builder.Len() == 0 || runes[i-1] != '_' {
				builder.WriteRune('_')
			}
			// 跳过所有连续的下划线
			for i+1 <= end && runes[i+1] == '_' {
				i++
			}
		} else {
			builder.WriteRune(current)
		}
	}

	return builder.String()
}
