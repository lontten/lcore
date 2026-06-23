package types

import "golang.org/x/exp/constraints"

// NewInt 将整数转为 *int，便于构造可选整数字段。
func NewInt[T constraints.Integer](i T) *int {
	t := int(i)
	return &t
}

// NewInt8 将整数转为 *int8，便于构造可选整数字段。
func NewInt8[T constraints.Integer](i T) *int8 {
	t := int8(i)
	return &t
}

// NewInt16 将整数转为 *int16，便于构造可选整数字段。
func NewInt16[T constraints.Integer](i T) *int16 {
	t := int16(i)
	return &t
}

// NewInt32 将整数转为 *int32，便于构造可选整数字段。
func NewInt32[T constraints.Integer](i T) *int32 {
	t := int32(i)
	return &t
}

// NewInt64 将整数转为 *int64，便于构造可选整数字段。
func NewInt64[T constraints.Integer](i T) *int64 {
	t := int64(i)
	return &t
}

// NewString 将字符串转为 *string，便于构造可选字符串字段。
func NewString(i string) *string {
	t := i
	return &t
}

// NewBool 将布尔值转为 *bool，便于构造可选布尔字段。
func NewBool(b bool) *bool {
	t := b
	return &t
}
