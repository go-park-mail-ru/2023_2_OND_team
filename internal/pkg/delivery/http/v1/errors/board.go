package errors

import (
	"errors"

	bCase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board"
)

// for backward compatibility
var (
	ErrEmptyTitle        = errors.New("empty or null board title has been provided")
	ErrEmptyPubOpt       = errors.New("null public option has been provided")
	ErrInvalidBoardTitle = errors.New("invalid or empty board title has been provided")
	ErrInvalidTagTitles  = errors.New("invalid tag titles have been provided")
	ErrInvalidUsername   = errors.New("invalid username has been provided")
)

var (
	WrappedErrors      = map[error]string{ErrInvalidTagTitles: "bad_Tagtitles"}
	ErrCodeCompability = map[error]string{
		ErrInvalidBoardTitle:     "bad_boardTitle",
		ErrEmptyTitle:            "empty_boardTitle",
		ErrEmptyPubOpt:           "bad_pubOpt",
		ErrInvalidUsername:       "bad_username",
		bCase.ErrInvalidUsername: "non_existingUser",
		bCase.ErrNoSuchBoard:     "no_board",
		bCase.ErrNoAccess:        "no_access",
		bCase.ErrNoPinOnBoard:    "no_pin",
	}
)

func GetErrCodeMessage(err error) (string, string) {
	var (
		code              string
		general, specific bool
	)

	code, general = generalErrCodeCompability[err]
	if general {
		return code, err.Error()
	}

	code, specific = ErrCodeCompability[err]
	if !specific {
		for wrappedErr, code_ := range WrappedErrors {
			if errors.Is(err, wrappedErr) {
				specific = true
				code = code_
			}
		}
	}
	if specific {
		return code, err.Error()
	}

	return ErrInternalError.Error(), generalErrCodeCompability[ErrInternalError]
}
