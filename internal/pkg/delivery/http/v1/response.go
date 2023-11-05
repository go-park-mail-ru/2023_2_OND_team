package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	InternalServerErrMessage = "internal server error occured"
	BadBodyMessage           = "can't parse body, JSON expected"
	BadQueryParamMessage     = "invalid query parameters have been provided"
)

var (
	BadBodyCode       = "bad_body"
	BadQueryParamCode = "bad_queryParam"
	InternalErrorCode = "internal_error"
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

func responseOk(w http.ResponseWriter, message string, body any) error {
	res := JsonResponse{
		Status:  "ok",
		Message: message,
		Body:    body,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("responseOk: %w", err)
	}
	w.WriteHeader(http.StatusOK)
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
