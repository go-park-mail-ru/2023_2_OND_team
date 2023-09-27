package crypto

import (
	"encoding/hex"
	"math/rand"
	"time"
	"unsafe"

	"golang.org/x/crypto/argon2"
)

func NewRandomStr(length int) (string, error) {
	rand.Seed(time.Now().UTC().UnixNano())
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
