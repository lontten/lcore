package lcutils

// nil 返回 零值
func NilToZero[T any](t *T) T {
	if t == nil {
		var zero T
		return zero
	}
	return *t
}
