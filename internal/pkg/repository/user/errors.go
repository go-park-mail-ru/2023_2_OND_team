package user

import errPkg "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/errors"

type ErrNonExistingUser struct{}

func (e *ErrNonExistingUser) Error() string {
	return "user doesn't exist"
}

func (e *ErrNonExistingUser) Type() errPkg.Type {
	return errPkg.ErrNotFound
}
