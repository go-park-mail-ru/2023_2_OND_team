package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	chi "github.com/go-chi/chi/v5"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	usecase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/pin"
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

	count, minID, maxID, err := FetchValidParamForLoadTape(r.URL)
	if err != nil {
		h.log.Info("parse url query params", log.F{"error", err.Error()})
		err = responseError(w, "bad_params",
			"expected parameters: count(positive integer: [1; 1000]), lastID(positive integer, the absence of this parameter is equal to the value 0)")
	} else {
		h.log.Sugar().Infof("param: count=%d, minID=%d, maxID=%d", count, minID, maxID)
		pins, minID, maxID := h.pinCase.SelectNewPins(r.Context(), count, minID, maxID)
		err = responseOk(w, "pins received are sorted by id", map[string]any{
			"pins":  pins,
			"minID": minID,
			"maxID": maxID,
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

	newPin := &entity.Pin{Author: &user.User{}}
	newPin.Author.ID = r.Context().Value(auth.KeyCurrentUserID).(int)

	tags := r.FormValue("tags")
	titles := strings.Split(tags, ",")
	newPin.Tags = make([]entity.Tag, 0, len(titles))
	for _, title := range titles {
		newPin.Tags = append(newPin.Tags, entity.Tag{Title: title})
	}

	newPin.SetTitle(r.FormValue("title"))

	newPin.SetDescription(r.FormValue("description"))

	public := r.FormValue("public")

	isPublic, err := strconv.ParseBool(public)
	if err != nil {
		responseError(w, "bad_body", "parameter public should have boolean value")
		return
	}
	newPin.Public = isPublic

	picture, mime, err := r.FormFile("picture")
	if err != nil {
		err = responseError(w, "bad_body", "unable to get an image from the request body")
		if err != nil {
			h.log.Error(err.Error())
		}
		return
	}
	defer picture.Close()

	newPin.Picture, err = h.imgCase.UploadImage("pins/", mime.Header.Get("Content-Type"), mime.Size, picture)
	if err != nil {
		err = responseError(w, "bad_body", "failed to upload the file received in the body")
		if err != nil {
			h.log.Error(err.Error())
		}
		return
	}

	err = h.pinCase.CreateNewPin(r.Context(), newPin)
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

func (h *HandlerHTTP) DeletePin(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request on delete new pin", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	userID := r.Context().Value(auth.KeyCurrentUserID).(int)

	pinIdStr := chi.URLParam(r, "pinID")
	pinID, err := strconv.ParseInt(pinIdStr, 10, 64)
	if err != nil {
		h.log.Error(err.Error())
		err = responseError(w, "parse_url", "internal error")
		if err != nil {
			h.log.Error(err.Error())
		}
		return
	}

	err = h.pinCase.DeletePinFromUser(r.Context(), int(pinID), userID)
	if err != nil {
		h.log.Error(err.Error())
		err = responseError(w, "pin_del", "internal error")
	} else {
		err = responseOk(w, "ok", nil)
	}
	if err != nil {
		h.log.Error(err.Error())
	}
}

func (h *HandlerHTTP) EditPin(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request on edit pin", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	userID := r.Context().Value(auth.KeyCurrentUserID).(int)

	pinIdStr := chi.URLParam(r, "pinID")
	pinID, err := strconv.ParseInt(pinIdStr, 10, 64)
	if err != nil {
		h.log.Error(err.Error())
		err = responseError(w, "parse_url", "internal error")
		if err != nil {
			h.log.Error(err.Error())
		}
		return
	}

	_, _ = userID, pinID

	pinUpdate := usecase.NewPinUpdateData()

	err = json.NewDecoder(r.Body).Decode(pinUpdate)
	defer r.Body.Close()
	if err != nil {
		h.log.Info(err.Error())
		err = responseError(w, "parse_body", "could not read the data to change")
		if err != nil {
			h.log.Error(err.Error())
		}
	}
	err = h.pinCase.EditPinByID(r.Context(), int(pinID), userID, pinUpdate)
	if err != nil {
		h.log.Error(err.Error())
		err = responseError(w, "edit_pin", "internal error")
	} else {
		err = responseOk(w, "pin data has been successfully changed", nil)
	}

	if err != nil {
		h.log.Error(err.Error())
	}
}

func (h *HandlerHTTP) ViewPin(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request on view pin", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	pinIdStr := chi.URLParam(r, "pinID")
	pinID, err := strconv.ParseInt(pinIdStr, 10, 64)
	if err != nil {
		h.log.Error(err.Error())
		err = responseError(w, "parse_url", "internal error")
		if err != nil {
			h.log.Error(err.Error())
		}
		return
	}

	userID, ok := r.Context().Value(auth.KeyCurrentUserID).(int)
	if !ok {
		userID = usecase.UserUnknown
	}
	pin, err := h.pinCase.ViewAnPin(r.Context(), int(pinID), userID)
	if err != nil {
		h.log.Error(err.Error())
		err = responseError(w, "edit_pin", "internal error")
	} else {
		err = responseOk(w, "pin was successfully received", pin)
	}
	if err != nil {
		h.log.Error(err.Error())
	}
}
