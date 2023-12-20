package structs

//go:generate easyjson response.go

//easyjson:json
type JsonResponse struct {
	Status  string      `json:"status" example:"ok"`
	Message string      `json:"message" example:"Response message"`
	Body    interface{} `json:"body" extensions:"x-omitempty"`
} // @name JsonResponse

//easyjson:json
type JsonErrResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Error description"`
	Code    string `json:"code"`
} // @name JsonErrResponse
