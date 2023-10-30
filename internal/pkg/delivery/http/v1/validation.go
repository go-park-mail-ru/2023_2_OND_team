package v1

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
)

var (
	ErrCountParameterMissing = errors.New("the count parameter is missing")
	ErrBadParams             = errors.New("bad params")
)

type errorFields []string

func (b *errorFields) Error() string {
	return strings.Join(*b, ",")
}

func (b *errorFields) addInvalidField(fieldName string) {
	*b = append(*b, fieldName)
}

func (b *errorFields) Err() error {
	if len(*b) == 0 {
		return nil
	}
	return b
}

func FetchValidParamForLoadTape(u *url.URL) (count int, lastID int, err error) {
	if param := u.Query().Get("count"); len(param) > 0 {
		c, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("fetch count param for load tape: %w", err)
		}
		count = int(c)
	} else {
		return 0, 0, ErrCountParameterMissing
	}
	if param := u.Query().Get("lastID"); len(param) > 0 {
		last, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("fetch lastID param for load tape: %w", err)
		}
		lastID = int(last)
	}
	if count <= 0 || count > 1000 || lastID < 0 {
		return 0, 0, ErrBadParams
	}
	return
}

func IsValidUserForRegistration(user *user.User) error {
	invalidFields := new(errorFields)

	if !isValidPassword(user.Password) {
		invalidFields.addInvalidField("password")
	}
	if !isValidEmail(user.Email) {
		invalidFields.addInvalidField("email")
	}
	if !isValidUsername(user.Username) {
		invalidFields.addInvalidField("username")
	}

	return invalidFields.Err()
}

func isValidUsername(username string) bool {
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

func isValidEmail(email string) bool {
	return valid.IsEmail(email) && len(email) <= 50
}

func isValidPassword(password string) bool {
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

func isValidName(name string) bool {
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

func isValidSurname(surname string) bool {
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

func isValidAboutMe(info string) bool {
	if len(info) > 500 {
		return false
	}
	for _, r := range info {
		if !(unicode.IsNumber(r) || unicode.IsLetter(r)) {
			return false
		}
	}
	return true
}
