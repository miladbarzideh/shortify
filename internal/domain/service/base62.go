package service

import (
	"strings"

	"github.com/pjebs/optimus-go"
)

const base62Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// This package utilizes Knuth's Hashing Algorithm to transform your internal ids into another number to hide it from the public.
var opt = optimus.New(1580030173, 59260789, 1163945558)

func Base62EncodeWithObfuscatedID(id uint64) string {
	obfuscatedID := opt.Encode(id)
	var encoding strings.Builder
	for obfuscatedID > 0 {
		encoding.WriteByte(base62Charset[obfuscatedID%62])
		obfuscatedID = obfuscatedID / 62
	}

	return encoding.String()
}
