package voyage

// P is a utility that converts a value to a pointer.
func P[T any](input T) *T {
	return &input
}
