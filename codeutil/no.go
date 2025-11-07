package codeutil

import (
	"crypto/rand"
	"time"
)

// NewTimedRandomID32
// 格式 yyyymmddHHMMSS + 18位随机数,有字母
func NewTimedRandomID32() string {
	return NewTimedRandomID(32)
}

// NewTimedRandomID
// 格式 yyyymmddHHMMSS + n位随机数,有字母
func NewTimedRandomID(length int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	timestamp := time.Now().Format("20060102150405") // 14字符

	randomLength := length - len(timestamp)
	if randomLength <= 0 {
		return timestamp[:length]
	}

	randomBytes := make([]byte, randomLength)
	if _, err := rand.Read(randomBytes); err != nil {
		panic("failed to generate random ID: " + err.Error())
	}

	for i := range randomBytes {
		randomBytes[i] = charset[randomBytes[i]%byte(len(charset))]
	}
	return timestamp + string(randomBytes)
}

// NewTimedRandomNumberID32
// 格式 yyyymmddHHMMSS + 18位随机数字
func NewTimedRandomNumberID32() string {
	return NewTimedRandomNumberID(32)
}

// NewTimedRandomNumberID
// 格式 yyyymmddHHMMSS + n位随机数字
func NewTimedRandomNumberID(length int) string {
	const charset = "0123456789"                     // 仅数字
	timestamp := time.Now().Format("20060102150405") // 14字符

	randomLength := length - len(timestamp)
	if randomLength <= 0 {
		return timestamp[:length]
	}

	randomBytes := make([]byte, randomLength)
	if _, err := rand.Read(randomBytes); err != nil {
		panic("failed to generate random ID: " + err.Error())
	}

	for i := range randomBytes {
		randomBytes[i] = charset[randomBytes[i]%byte(len(charset))]
	}
	return timestamp + string(randomBytes)
}
