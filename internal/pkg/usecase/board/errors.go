package board

import "errors"

var (
	ErrInvalidUsername   = errors.New("invalid username has been provided or username doesn't exist")
	ErrNoSuchBoard       = errors.New("board is not accessable or doesn't exist")
	ErrInvalidBoardTitle = errors.New("invalid or empty board title has been provided")
	ErrInvalidTagTitles  = errors.New("invalid tag titles have been provided")
	ErrInvalidUserID     = errors.New("invalid user id has been provided")
)
