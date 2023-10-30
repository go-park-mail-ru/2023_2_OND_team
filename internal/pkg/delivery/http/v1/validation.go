package v1

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	valid "github.com/go-park-mail-ru/2023_2_OND_team/pkg/validation"
)

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
	invalidFields := new(valid.ErrorFields)

	if !valid.IsValidPassword(user.Password) {
		invalidFields.AddInvalidField("password")
	}
	if !valid.IsValidEmail(user.Email) {
		invalidFields.AddInvalidField("email")
	}
	if !valid.IsValidUsername(user.Username) {
		invalidFields.AddInvalidField("username")
	}

	return invalidFields.Err()
}
