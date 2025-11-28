package cmd

func toPtr[T any](x T) *T {
	return &x
}
