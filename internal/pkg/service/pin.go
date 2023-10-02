package service

import (
	"errors"
	"net/http"

	_ "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"

	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

var ErrCountParameterMissing = errors.New("the count parameter is missing")
var ErrBadParams = errors.New("bad params")

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
	s.log.Info("request on get pins", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	count, lastID, err := FetchValidParamForLoadTape(r.URL)
	if err != nil {
		s.log.Info("parse url query params", log.F{"error", err.Error()})
		err = responseError(w, "bad_params",
			"expected parameters: count(positive integer: [1; 1000]), lastID(positive integer, the absence of this parameter is equal to the value 0)")
	} else {
		s.log.Sugar().Infof("param: count=%d, lastID=%d", count, lastID)
		pins, last := s.pinCase.SelectNewPins(r.Context(), count, lastID)
		err = responseOk(w, "pins received are sorted by id", map[string]any{
			"pins":   pins,
			"lastID": last,
		})
	}
	if err != nil {
		s.log.Error(err.Error())
	}
}
