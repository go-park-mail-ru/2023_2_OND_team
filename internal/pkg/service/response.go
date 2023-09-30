package service

import (
	"net/http"
)

type Empty struct{} // @name Empty

type JsonResponse struct {
	Status  string      `json:"status" example:"ok"`
	Message string      `json:"message" example:"Response message"`
	Body    interface{} `json:"body"`
} // @name JsonResponse

type JsonErrResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Error description"`
	Code    int    `json:"code,string"`
} // @name JsonErrResponse

func SetContentTypeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
