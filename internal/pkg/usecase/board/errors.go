package board

import "errors"

var (
	ErrInvalidUsername = errors.New("username doesn't exist")
	ErrNoSuchBoard     = errors.New("board is not accessable or doesn't exist")
	ErrNoPinOnBoard    = errors.New("no such pin on board")
	ErrInvalidUserID   = errors.New("invalid user id has been provided")
	ErrNoAccess        = errors.New("no access for this action")
)
