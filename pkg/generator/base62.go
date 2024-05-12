package generator

import (
	"math/rand"
	"regexp"
	"time"
)

const base62Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var base62Regex = regexp.MustCompile("^[a-zA-Z0-9]+$")

type Generator interface {
	GenerateShortURLCode(length int) string
}

type generator struct {
	rand *rand.Rand
}

func NewGenerator() Generator {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &generator{
		rand: r,
	}
}

func (g *generator) GenerateShortURLCode(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		charIndex := g.rand.Intn(len(base62Charset))
		result[i] = base62Charset[charIndex]
	}

	return string(result)
}

func IsValidBase62(s string) bool {
	return base62Regex.MatchString(s)
}
