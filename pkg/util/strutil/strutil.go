package strutil

import (
	"math/rand"
	"time"

	"github.com/teris-io/shortid"
)

const (
	Number       = "0123456789"
	Alphabet     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Alphanumeric = Alphabet + Number
)

var (
	numberRune       = []rune(Number)
	alphabetRune     = []rune(Alphabet)
	alphanumericRune = []rune(Alphanumeric)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// ShortID generate short id. Use RandomAlphanumeric(10) instead When short id generate fail.
func ShortID() string {
	if id, err := shortid.Generate(); err == nil {
		return id
	}

	return RandomAlphanumeric(10)
}

// RandomNumber return a random n length string of number.
func RandomNumber(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = numberRune[rand.Intn(len(numberRune))]
	}

	return string(b)
}

// RandomAlphabet return a random n length string of alphabet.
func RandomAlphabet(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = alphabetRune[rand.Intn(len(alphabetRune))]
	}

	return string(b)
}

// RandomAlphanumeric return a random n length string of alphanumeric.
func RandomAlphanumeric(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = alphanumericRune[rand.Intn(len(alphanumericRune))]
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

func Compact(array []string) []string {
	var newArr []string

	for _, arr := range array {
		if arr != "" {
			newArr = append(newArr, arr)
		}
	}

	return newArr
}
