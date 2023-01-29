package utils

import "fmt"

func GetAddress[T any](s T) *T {
	return &s
}

func URLFromToken(token string) string {
	return fmt.Sprintf("https://beneburg.edyapups.ru/login/%s", token)
}
