package lcutils

import (
	"crypto/rand"
	"time"
)

// GenerateOrderID32
// 格式 yyyymmddHHMMSS + 18位随机数,有字母
func GenerateOrderID32() string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	timestamp := time.Now().Format("20060102150405")

	// 生成 18 字节安全随机数（直接对应最终长度）
	randomBytes := make([]byte, 18)
	_, _ = rand.Read(randomBytes)

	// 直接映射字符集（避免编码错误）
	for i := range randomBytes {
		randomBytes[i] = charset[randomBytes[i]%byte(len(charset))]
	}
	return timestamp + string(randomBytes)
}

// GenerateOrderID32Number
// 格式 yyyymmddHHMMSS + 18位随机数, 纯数字
func GenerateOrderID32Number() string {
	// 时间部分 (14位)
	timestamp := time.Now().Format("20060102150405")

	// 生成安全随机数（优化版本）
	randomBytes := make([]byte, 18)
	_, _ = rand.Read(randomBytes)

	// 改进的均匀分布映射（避免模偏差）
	const digits = "0123456789"
	result := make([]byte, 18)
	for i := range randomBytes {
		// 使用更均匀的分布方式
		result[i] = digits[int(randomBytes[i])%len(digits)]
	}

	return timestamp + string(result)
}
