package strutil

import (
	"math/rand"
	"time"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandString return a random n length string
func RandString(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// First return first element of string array.
func First(array []string) string {
	return Nth(array, 0)
}

// Last return last element of string array.
func Last(array []string) string {
	return Nth(array, len(array)-1)
}

// Nth return the nth element. if index beyond the length, return "".
func Nth(array []string, idx int) string {
	if idx < len(array) {
		return array[idx]
	}

	return ""
}
