package codeutil

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func GetSHA256HashCode(message []byte, salt string) string {
	//计算哈希值,返回一个长度为32的数组
	bytes2 := sha256.Sum256([]byte(string(message) + salt))
	//将数组转换成切片,转换成16进制,返回字符串
	hashcode := hex.EncodeToString(bytes2[:])
	return hashcode
}

// base64 加密
func Base64Encode(str string) string {
	input := []byte(str)
	return base64.StdEncoding.EncodeToString(input)
}

// base64 解密
func Base64Decode(str string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(str)
	return string(decodeBytes), err
}
