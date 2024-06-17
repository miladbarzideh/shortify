package generator

import (
	"math/rand"
	"regexp"
	"time"
)

const base62Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var base62Regex = regexp.MustCompile("^[a-zA-Z0-9]+$")

type RandomGenerator struct {
	rand   *rand.Rand
	length int
}

func NewGenerator(length int) *RandomGenerator {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &RandomGenerator{
		rand:   r,
		length: length,
	}
}

func (g *RandomGenerator) GenerateShortURLCode() string {
	result := make([]byte, g.length)
	for i := 0; i < g.length; i++ {
		charIndex := g.rand.Intn(len(base62Charset))
		result[i] = base62Charset[charIndex]
	}

	return string(result)
}

func IsValidBase62(s string) bool {
	return base62Regex.MatchString(s)
}

func (g *RandomGenerator) SetLength(length int) {
	g.length = length
}
