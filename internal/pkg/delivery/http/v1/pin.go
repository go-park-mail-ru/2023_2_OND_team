package v1

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
//	@Accept			json
//	@Produce		json
//	@Param			lastID	path		string	false	"ID of the pin that will be just before the first pin in the requested collection, 0 by default"	example(2)
//
// @Param			count	path		string	true	"Pins quantity after last pin specified in lastID"													example(5)
// @Success		200		{object}	JsonResponse{body=[]Pin}
// @Failure		400		{object}	JsonErrResponse
// @Failure		404		{object}	JsonErrResponse
// @Failure		500		{object}	JsonErrResponse
// @Router			/api/v1/pin [get]
func (h *HandlerHTTP) GetPins(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request on get pins", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	count, lastID, err := FetchValidParamForLoadTape(r.URL)
	if err != nil {
		h.log.Info("parse url query params", log.F{"error", err.Error()})
		err = responseError(w, "bad_params",
			"expected parameters: count(positive integer: [1; 1000]), lastID(positive integer, the absence of this parameter is equal to the value 0)")
	} else {
		h.log.Sugar().Infof("param: count=%d, lastID=%d", count, lastID)
		pins, last := h.pinCase.SelectNewPins(r.Context(), count, lastID)
		err = responseOk(w, "pins received are sorted by id", map[string]any{
			"pins":   pins,
			"lastID": last,
		})
	}
	if err != nil {
		h.log.Error(err.Error())
	}
}