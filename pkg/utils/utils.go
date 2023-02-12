package utils

func GetAddress[T any](s T) *T {
	return &s
}
