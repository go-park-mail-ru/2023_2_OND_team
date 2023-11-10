package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

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
