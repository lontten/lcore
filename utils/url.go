package lcutils

import (
	"regexp"
	"strings"
)

func SanitizeURL(url string) string {
	url = strings.ReplaceAll(url, ":", "_")
	// 步骤1: 定义正则表达式规则
	// 允许的字符范围：
	// - 中日韩文字符：\p{Han}（中文/日文汉字）、\p{Hiragana}（平假名）、\p{Katakana}（片假名）、\p{Hangul}（韩文字母）
	// - 拉丁字母、数字、下划线、点号、短横线
	illegalPattern := `[^\p{Han}\p{Hiragana}\p{Katakana}\p{Hangul}a-zA-Z0-9_.-]+`
	// 预编译正则表达式（提升性能）
	reg := regexp.MustCompile(illegalPattern)
	// 步骤2: 替换非法字符为短横线
	sanitized := reg.ReplaceAllString(url, "-")
	return sanitized
}
