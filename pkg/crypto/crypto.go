package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"unsafe"

	"golang.org/x/crypto/argon2"
)

var ErrNegativeLen = errors.New("the length cannot be negative")

func NewRandomString(length int) (string, error) {
	if length < 0 {
		return "", ErrNegativeLen
	}
	b := make([]byte, hex.DecodedLen(length))
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	res := make([]byte, length)
	hex.Encode(res, b)
	return *(*string)(unsafe.Pointer(&res)), nil
}

func PasswordHash(password, salt string, length int) string {
	if length <= 0 {
		return ""
	}
	passwordHash := make([]byte, length)
	hex.Encode(passwordHash, argon2.IDKey([]byte(password), []byte(salt), 3, 32*1024, 4,
		uint32(hex.DecodedLen(length))))
	return *(*string)(unsafe.Pointer(&passwordHash))
}
