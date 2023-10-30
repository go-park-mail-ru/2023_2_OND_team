package validation

import (
	"strings"
	"unicode"

	valid "github.com/asaskevich/govalidator"
)

type ErrorFields []string

func (b *ErrorFields) Error() string {
	return strings.Join(*b, ",")
}

func (b *ErrorFields) AddInvalidField(fieldName string) {
	*b = append(*b, fieldName)
}

func (b *ErrorFields) Err() error {
	if len(*b) == 0 {
		return nil
	}
	return b
}

func IsValidUsername(username string) bool {
	if len(username) < 4 || len(username) > 50 {
		return false
	}
	for _, r := range username {
		if !(unicode.IsNumber(r) || unicode.IsLetter(r)) {
			return false
		}
	}
	return true
}

func IsValidEmail(email string) bool {
	return valid.IsEmail(email) && len(email) <= 50
}

func IsValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 50 {
		return false
	}
	for _, r := range password {
		if !(unicode.IsNumber(r) || unicode.IsSymbol(r) || unicode.IsPunct(r) || unicode.IsLetter(r)) {
			return false
		}
	}
	return true
}

func IsValidName(name string) bool {
	if len(name) > 50 {
		return false
	}
	for _, r := range name {
		if !(unicode.IsNumber(r) || unicode.IsLetter(r)) {
			return false
		}
	}
	return true
}

func IsValidSurname(surname string) bool {
	if len(surname) > 50 {
		return false
	}
	for _, r := range surname {
		if !(unicode.IsNumber(r) || unicode.IsLetter(r)) {
			return false
		}
	}
	return true
}
