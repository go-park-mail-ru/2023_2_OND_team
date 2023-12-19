package v1

import (
	"fmt"
	"net/http"

	errHTTP "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1/errors"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1/structs"
	"github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
	"github.com/mailru/easyjson"
)

func SetContentTypeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func responseOk(statusCode int, w http.ResponseWriter, message string, body any) error {
	res := structs.JsonResponse{
		Status:  "ok",
		Message: message,
		Body:    body,
	}
	resBytes, err := easyjson.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("responseOk: %w", err)
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(resBytes)
	return err
}

func responseError(w http.ResponseWriter, code, message string) error {
	res := structs.JsonErrResponse{
		Status:  "error",
		Message: message,
		Code:    code,
	}
	resBytes, err := easyjson.Marshal(res)
	if err != nil {
		return fmt.Errorf("responseError: %w", err)
	}
	_, err = w.Write(resBytes)
	return err
}

func (h *HandlerHTTP) responseErr(w http.ResponseWriter, r *http.Request, err error) error {
	log := logger.GetLoggerFromCtx(r.Context())

	code, status := errHTTP.GetCodeStatusHttp(err)
	var msg string
	if status == http.StatusInternalServerError {
		log.Warnf("unexpected application error: %s", err.Error())
		msg = "internal error occured"
	} else {
		msg = err.Error()
	}

	res := structs.JsonErrResponse{
		Status:  "error",
		Message: msg,
		Code:    code,
	}
	resBytes, err := easyjson.Marshal(res)
	if err != nil {
		return fmt.Errorf("responseError: %w", err)
	}
	w.WriteHeader(status)
	_, err = w.Write(resBytes)
	return err
}
