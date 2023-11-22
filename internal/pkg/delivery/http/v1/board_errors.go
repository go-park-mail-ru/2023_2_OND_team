package v1

import (
	"errors"

	bCase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/board"
)

var (
	ErrEmptyTitle        = errors.New("empty or null board title has been provided")
	ErrEmptyPubOpt       = errors.New("null public option has been provided")
	ErrInvalidBoardTitle = errors.New("invalid or empty board title has been provided")
	ErrInvalidTagTitles  = errors.New("invalid tag titles have been provided")
	ErrInvalidUsername   = errors.New("invalid username has been provided")
)

var (
	wrappedErrors      = map[error]string{ErrInvalidTagTitles: "bad_Tagtitles"}
	errCodeCompability = map[error]string{
		ErrInvalidBoardTitle:     "bad_boardTitle",
		ErrEmptyTitle:            "empty_boardTitle",
		ErrEmptyPubOpt:           "bad_pubOpt",
		ErrInvalidUsername:       "bad_username",
		bCase.ErrInvalidUsername: "non_existingUser",
		bCase.ErrNoSuchBoard:     "no_board",
		bCase.ErrNoAccess:        "no_access",
	}
)

func getErrCodeMessage(err error) (string, string) {
	var (
		code              string
		general, specific bool
	)

	code, general = generalErrCodeCompability[err]
	if general {
		return code, err.Error()
	}

	code, specific = errCodeCompability[err]
	if !specific {
		for wrappedErr, code_ := range wrappedErrors {
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
