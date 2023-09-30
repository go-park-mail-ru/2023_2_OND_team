package service

import (
	"fmt"
	"net/http"

	_ "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
)

// GetPins godoc
//
//	@Description	Get pin collection
//	@Tags			Pin
//	@Produce		json
//	@Success		200	{object}	JsonResponse{body=[]Pin}
//	@Failure		400	{object}	JsonErrResponse
//	@Failure		404	{object}	JsonErrResponse
//	@Failure		500	{object}	JsonErrResponse
//	@Router			/api/v1/pin [get]
func (s *Service) GetPins(w http.ResponseWriter, r *http.Request) {
	s.log.Info("it worked GetPins")
	fmt.Fprintf(w, "{\"status\": \"ok\", \"path\": \"%s\", \"method\": \"%s\"}\n", r.URL.Path, r.Method)
}

// GetPinByID godoc
//
//	@Description	Get concrete pin by id
//	@Tags			Pin
//	@Produce		json
//	@Param			pinId	path		int	true	"Id of the pin"
//	@Success		200		{object}	JsonResponse{body=Pin}
//	@Failure		400		{object}	JsonErrResponse
//	@Failure		404		{object}	JsonErrResponse
//	@Failure		500		{object}	JsonErrResponse
//	@Router			/api/v1/pin/{pinId} [get]
func (s *Service) GetPinByID(w http.ResponseWriter, r *http.Request) {
	s.log.Info("it worked GetPinByID")
	fmt.Fprintf(w, "{\"status\": \"ok\", \"path\": \"%s\", \"method\": \"%s\"}\n", r.URL.Path, r.Method)
}
