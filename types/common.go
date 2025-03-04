package types

import "golang.org/x/exp/constraints"

func NewInt[T constraints.Integer](i T) *int {
	var t = int(i)
	return &t
}

func NewInt8[T constraints.Integer](i T) *int8 {
	var t = int8(i)
	return &t
}

func NewInt16[T constraints.Integer](i T) *int16 {
	var t = int16(i)
	return &t
}

func NewInt32[T constraints.Integer](i T) *int32 {
	var t = int32(i)
	return &t
}

func NewInt64[T constraints.Integer](i T) *int64 {
	var t = int64(i)
	return &t
}

func NewString(i string) *string {
	var t = i
	return &t
}

func NewBool(b bool) *bool {
	var t = b
	return &t
}
