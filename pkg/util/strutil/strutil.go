package strutil

import (
	"math/rand"
	"reflect"
	"strings"
	"time"
	"unicode"
	"unsafe"

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

// Compact remove empty string element.
func Compact(array []string) []string {
	var newArr []string

	for _, arr := range array {
		if arr != "" {
			newArr = append(newArr, arr)
		}
	}

	return newArr
}

// HasPrefix check any element in arr has prefix
func HasPrefix(arr []string, prefix string) (b bool) {
	return ContainsBy(arr, func(str string) bool {
		return strings.HasPrefix(str, prefix)
	})
}

// HasSuffix check any element in arr has suffix
func HasSuffix(arr []string, suffix string) (b bool) {
	return ContainsBy(arr, func(str string) bool {
		return strings.HasSuffix(str, suffix)
	})
}

// Contains check any element in arr equal to flag
func Contains(arr []string, flag string) (b bool) {
	return ContainsBy(arr, func(str string) bool {
		return str == flag
	})
}

// Contained check any element in arr contained by flag
func Contained(arr []string, flag string) (b bool) {
	return ContainsBy(arr, func(str string) bool {
		return strings.Contains(flag, str)
	})
}

// ContainsBy check any element in arr by fn
func ContainsBy(arr []string, fn func(str string) bool) (b bool) {
	for _, r := range arr {
		if fn(r) {
			b = true
			break
		}
	}

	return
}

// ContainsChineseChar check str contains chinese char
func ContainsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}

// String2Bytes more effect but not safe
func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// Bytes2String more effect but not safe
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
