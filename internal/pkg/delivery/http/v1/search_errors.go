package v1

import errPkg "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/errors"

type ErrNoData struct{}

func (e *ErrNoData) Error() string {
	return "Can't find any user/board/pin"
}

func (e *ErrNoData) Type() errPkg.Type {
	return errPkg.ErrNotFound
}

type ErrInvalidTemplate struct{}

func (e *ErrInvalidTemplate) Error() string {
	return "Invalid template has been provided"
}

func (e *ErrInvalidTemplate) Type() errPkg.Type {
	return errPkg.ErrInvalidInput
}
