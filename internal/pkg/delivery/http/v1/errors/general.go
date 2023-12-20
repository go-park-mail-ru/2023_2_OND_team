package errors

import (
	"errors"
	"fmt"
	"net/http"

	errPkg "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/errors"
)

// for backward compatibility - begin
var (
	ErrBadBody        = errors.New("can't parse body, JSON with correct data types is expected")
	ErrBadUrlParam    = errors.New("bad URL param has been provided")
	ErrBadQueryParam  = errors.New("invalid query parameters have been provided")
	ErrInternalError  = errors.New("internal server error occured")
	ErrBadContentType = errors.New("application/json is expected")
)

var (
	generalErrCodeCompability = map[error]string{
		ErrBadBody:        "bad_body",
		ErrBadQueryParam:  "bad_queryParams",
		ErrInternalError:  "internal_error",
		ErrBadContentType: "bad_contentType",
		ErrBadUrlParam:    "bad_urlParam",
	}
)

// for backward compatibility - end

type ErrInvalidBody struct{}

func (e *ErrInvalidBody) Error() string {
	return "invalid body"
}

func (e *ErrInvalidBody) Type() errPkg.Type {
	return errPkg.ErrInvalidInput
}

type ErrInvalidQueryParam struct {
	Params map[string]string
}

func (e *ErrInvalidQueryParam) Error() string {
	return fmt.Sprintf("invalid query params: %v", e.Params)
}

func (e *ErrInvalidQueryParam) Type() errPkg.Type {
	return errPkg.ErrInvalidInput
}

type ErrInvalidContentType struct {
	PreferredType string
}

func (e *ErrInvalidContentType) Error() string {
	return fmt.Sprintf("invalid content type, should be %s", e.PreferredType)
}

func (e *ErrInvalidContentType) Type() errPkg.Type {
	return errPkg.ErrInvalidInput
}

type ErrInvalidUrlParams struct {
	Params map[string]string
}

func (e *ErrInvalidUrlParams) Error() string {
	return fmt.Sprintf("invalid URL params: %v", e.Params)
}

func (e *ErrInvalidUrlParams) Type() errPkg.Type {
	return errPkg.ErrInvalidInput
}

type ErrMissingBodyParams struct {
	Params []string
}

func (e *ErrMissingBodyParams) Error() string {
	return fmt.Sprintf("missing body params: %v", e.Params)
}

func (e *ErrMissingBodyParams) Type() errPkg.Type {
	return errPkg.ErrInvalidInput
}

func GetCodeStatusHttp(err error) (ErrCode string, httpStatus int) {

	var declaredErr errPkg.DeclaredError
	if errors.As(err, &declaredErr) {
		switch declaredErr.Type() {
		case errPkg.ErrInvalidInput:
			return "bad_input", http.StatusBadRequest
		case errPkg.ErrNotFound:
			return "not_found", http.StatusNotFound
		case errPkg.ErrAlreadyExists:
			return "already_exists", http.StatusConflict
		case errPkg.ErrNoAuth:
			return "no_auth", http.StatusUnauthorized
		case errPkg.ErrNoAccess:
			return "no_access", http.StatusForbidden
		case errPkg.ErrTimeout:
			return "timeout", http.StatusRequestTimeout
		}
	}

	return "internal_error", http.StatusInternalServerError
}
