package generator

import (
	"crypto/rand"
	"math/big"
	"regexp"
)

const base62Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var base62Regex = regexp.MustCompile("^[a-zA-Z0-9]+$")

func GenerateShortURLCode(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		charIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(base62Charset))))
		result[i] = base62Charset[charIndex.Int64()]
	}
	return string(result)
}

func IsValidBase62(s string) bool {
	return base62Regex.MatchString(s)
}
