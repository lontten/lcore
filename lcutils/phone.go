package lcutils

import (
	"regexp"
)

func CheckPhoneAll(phoneNumber string) bool {
	// 手机号
	pattern := "^1[3-9][0-9]{9}$"
	regex := regexp.MustCompile(pattern)
	matched := regex.MatchString(phoneNumber)
	if matched {
		return true
	}

	// 固话 区号（2-3位） + 连字符 + 本地号码（7-8位） + 连字符 + 分机号（2-4位）
	pattern = `^0\d{2,3}-\d{7,8}(-\d{2,4})?$`
	regex = regexp.MustCompile(pattern)
	matched = regex.MatchString(phoneNumber)
	if matched {
		return true
	}
	return false
}

func CheckPhone(phoneNumber string) bool {
	// 手机号
	pattern := "^1[3-9][0-9]{9}$"
	regex := regexp.MustCompile(pattern)
	matched := regex.MatchString(phoneNumber)
	if matched {
		return true
	}
	return false
}

func CheckLandline(phoneNumber string) bool {
	// 固话 区号（2-3位） + 连字符 + 本地号码（7-8位） + 连字符 + 分机号（2-4位）
	pattern := `^0\d{2,3}-\d{7,8}(-\d{2,4})?$`
	regex := regexp.MustCompile(pattern)
	matched := regex.MatchString(phoneNumber)
	if matched {
		return true
	}
	return false
}
