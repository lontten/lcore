package utils

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

	// 固话
	pattern = `^0\d{2,3}-\d{7,8}$`
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
	// 固话
	pattern := `^0\d{2,3}-\d{7,8}$`
	regex := regexp.MustCompile(pattern)
	matched := regex.MatchString(phoneNumber)
	if matched {
		return true
	}
	return false
}
