package v1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

const MaxMemoryParseFormData = 10 * 1 << 20

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

func (h *HandlerHTTP) CreateNewPin(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request on create new pin", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		err := responseError(w, "bad_request", "the request body should be multipart/form-data")
		if err != nil {
			h.log.Error(err.Error())
		}
		return
	}

	err := r.ParseMultipartForm(MaxMemoryParseFormData)
	if err != nil {
		err = responseError(w, "bad_body", "failed to read request body")
		if err != nil {
			h.log.Error(err.Error())
		}
		return
	}
	defer r.Body.Close()

	newPin := &pin.Pin{}
	newPin.AuthorID = r.Context().Value(auth.KeyCurrentUserID).(int)

	tags := r.FormValue("tags")
	_ = tags
	newPin.Tags = []pin.Tag{pin.Tag{Title: "good"}, pin.Tag{Title: "aaa"}, pin.Tag{Title: "bbb"}}

	newPin.Title = r.FormValue("title")

	newPin.Description = r.FormValue("description")

	public := r.FormValue("public")

	isPublic, err := strconv.ParseBool(public)
	if err != nil {
		responseError(w, "bad_body", "parameter public should have boolean value")
		return
	}
	newPin.Public = isPublic

	picture, mime, err := r.FormFile("picture")
	if err != nil {
		responseError(w, "bad_body", "unable to get an image from the request body")
		return
	}
	defer picture.Close()

	err = h.pinCase.CreateNewPin(r.Context(), newPin, picture, mime.Header.Get("Content-Type"))
	if err != nil {
		h.log.Error(err.Error())
		err = responseError(w, "add_pin", "failed to create pin")
	} else {
		err = responseOk(w, "pin successfully created", nil)
	}
	if err != nil {
		h.log.Error(err.Error())
	}
}
