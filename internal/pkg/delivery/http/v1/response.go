package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	errPkg "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/errors"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
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
	params map[string]string
}

func (e *ErrInvalidQueryParam) Error() string {
	return fmt.Sprintf("invalid query params: %v", e.params)
}

func (e *ErrInvalidQueryParam) Type() errPkg.Type {
	return errPkg.ErrInvalidInput
}

type ErrInvalidContentType struct{}

func (e *ErrInvalidContentType) Error() string {
	return "invalid content type"
}

func (e *ErrInvalidContentType) Type() errPkg.Type {
	return errPkg.ErrInvalidInput
}

type ErrInvalidUrlParams struct {
	params map[string]string
}

func (e *ErrInvalidUrlParams) Error() string {
	return fmt.Sprintf("invalid URL params: %v", e.params)
}

func (e *ErrInvalidUrlParams) Type() errPkg.Type {
	return errPkg.ErrInvalidInput
}

type ErrMissingBodyParams struct {
	params []string
}

func (e *ErrMissingBodyParams) Error() string {
	return fmt.Sprintf("missing body params: %v", e.params)
}

func (e *ErrMissingBodyParams) Type() errPkg.Type {
	return errPkg.ErrInvalidInput
}

func getCodeStatusHttp(err error) (ErrCode string, httpStatus int) {

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

type JsonResponse struct {
	Status  string      `json:"status" example:"ok"`
	Message string      `json:"message" example:"Response message"`
	Body    interface{} `json:"body" extensions:"x-omitempty"`
} // @name JsonResponse

type JsonErrResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Error description"`
	Code    string `json:"code"`
} // @name JsonErrResponse

func SetContentTypeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func responseOk(statusCode int, w http.ResponseWriter, message string, body any) error {
	res := JsonResponse{
		Status:  "ok",
		Message: message,
		Body:    body,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("responseOk: %w", err)
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(resBytes)
	return err
}

func responseError(w http.ResponseWriter, code, message string) error {
	res := JsonErrResponse{
		Status:  "error",
		Message: message,
		Code:    code,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("responseError: %w", err)
	}
	_, err = w.Write(resBytes)
	return err
}

func (h *HandlerHTTP) responseErr(w http.ResponseWriter, r *http.Request, err error) error {
	log := logger.GetLoggerFromCtx(r.Context())

	code, status := getCodeStatusHttp(err)
	var msg string
	if status == http.StatusInternalServerError {
		log.Warnf("unexpected error on the delivery http: %s\n", err.Error())
		err := &errPkg.InternalError{}
		msg = err.Error()
	} else {
		msg = err.Error()
	}

	res := JsonErrResponse{
		Status:  "error",
		Message: msg,
		Code:    code,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("responseError: %w", err)
	}
	w.WriteHeader(status)
	_, err = w.Write(resBytes)
	return err
}
