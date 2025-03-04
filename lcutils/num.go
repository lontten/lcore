package lcutils

import (
	"fmt"
	"strconv"
)
import "golang.org/x/exp/constraints"

// Num2Str 将任意整数类型的数值转换为字符串
func Num2Str[T constraints.Integer](num T) string {
	return fmt.Sprintf("%d", num)
}

// Str2Num 将字符串转换为int64数值，支持任意进制解析
func Str2Num(str string) (int64, error) {
	// 自动识别进制（0b/0o/0x前缀，否则十进制）
	return strconv.ParseInt(str, 0, 64)
}
func Str2NumMust(str string) int64 {
	num, err := Str2Num(str)
	if err != nil {
		panic(err)
	}
	return num
}
