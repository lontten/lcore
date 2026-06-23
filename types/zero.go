package types

// NilToZero 当 t 为 nil 时返回 T 的零值，否则返回 *t。
func NilToZero[T any](t *T) T {
	if t == nil {
		var zero T
		return zero
	}
	return *t
}

// NilToZeroP 当 t 为 nil 时返回指向 T 零值的指针，否则返回 t。
func NilToZeroP[T any](t *T) *T {
	if t == nil {
		var zero T
		return &zero
	}
	return t
}
