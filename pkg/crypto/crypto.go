package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"unsafe"

	"golang.org/x/crypto/argon2"
)

func NewRandomString(length int) (string, error) {
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
	passwordHash := make([]byte, length)
	hex.Encode(passwordHash, argon2.IDKey([]byte(password), []byte(salt), 3, 32*1024, 4,
		uint32(hex.DecodedLen(length))))
	return *(*string)(unsafe.Pointer(&passwordHash))
}
