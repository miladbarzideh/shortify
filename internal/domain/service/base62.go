package service

import (
	"strings"
)

const base62Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Base62Encode(n uint64) string {
	var encoding strings.Builder
	for n > 0 {
		encoding.WriteByte(base62Charset[n%62])
		n = n / 62
	}

	return encoding.String()
}
