package types

// nil 返回 零值
func NilToZero[T any](t *T) T {
	if t == nil {
		var zero T
		return zero
	}
	return *t
}

// 指针为nil时，返回零值指针
func NilToZeroP[T any](t *T) *T {
	if t == nil {
		var zero T
		return &zero
	}
	return t
}
