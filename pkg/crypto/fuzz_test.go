package crypto

import (
	"testing"
)

func FuzzNewRandomString(f *testing.F) {
	f.Add(5)
	f.Add(0)
	f.Add(-2)

	f.Fuzz(func(t *testing.T, length int) {
		length %= 10000
		randStr, err := NewRandomString(length)
		if length < 0 && (randStr != "" || err != ErrNegativeLen) || length >= 0 && (len(randStr) != length || err != nil) {
			t.Fatalf("NewRandomString(%d) retured %s, %v, lenght returned string equal %d, but except %d",
				length, randStr, err, len(randStr), length)
		}
	})
}

func FuzzPasswordHash(f *testing.F) {
	f.Add("password", "salt", 5)
	f.Add("a", "apple", 0)
	f.Add("", "", -1)

	f.Fuzz(func(t *testing.T, password, salt string, length int) {
		length %= 10000
		passHash := PasswordHash(password, salt, length)
		if length < 0 && passHash != "" || length >= 0 && len(passHash) != length {
			t.Fatalf("PasswordHash(%s, %s, %d) retured %s, lenght returned string equal %d, but except %d",
				password, salt, length, passHash, len(passHash), length)
		}
	})
}
