package v1

import (
	"net/http"
	"strconv"
	"strings"

	chi "github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"

	entity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	img "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/image"
	usecase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/pin"
)

const MaxMemoryParseFormData = 12 * 1 << 20

func (h *HandlerHTTP) CreateNewPin(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		err := responseError(w, "bad_request", "the request body should be multipart/form-data")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	err := r.ParseMultipartForm(MaxMemoryParseFormData)
	if err != nil {
		err = responseError(w, "bad_body", "failed to read request body")
		if err != nil {
			logger.Error(err.Error())
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
			logger.Error(err.Error())
		}
		return
	}
	defer picture.Close()

	err = h.pinCase.CreateNewPin(r.Context(), newPin, mime.Header.Get("Content-Type"), mime.Size, picture)
	if err != nil {
		logger.Error(err.Error())
		if err == img.ErrExplicitImage {
			err = responseError(w, "explicit_pin", err.Error())
		} else {
			err = responseError(w, "add_pin", "failed to create pin")
		}
	} else {
		err = responseOk(http.StatusCreated, w, "pin successfully created", nil)
	}
	if err != nil {
		logger.Error(err.Error())
	}
}

func (h *HandlerHTTP) DeletePin(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	userID := r.Context().Value(auth.KeyCurrentUserID).(int)

	pinIdStr := chi.URLParam(r, "pinID")
	pinID, err := strconv.ParseInt(pinIdStr, 10, 64)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "parse_url", "internal error")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	err = h.pinCase.DeletePinFromUser(r.Context(), int(pinID), userID)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "pin_del", "internal error")
	} else {
		err = responseOk(http.StatusOK, w, "ok", nil)
	}
	if err != nil {
		logger.Error(err.Error())
	}
}

func (h *HandlerHTTP) EditPin(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	userID := r.Context().Value(auth.KeyCurrentUserID).(int)

	pinIdStr := chi.URLParam(r, "pinID")
	pinID, err := strconv.ParseInt(pinIdStr, 10, 64)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "parse_url", "internal error")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	_, _ = userID, pinID

	pinUpdate := &usecase.PinUpdateData{}
	err = easyjson.UnmarshalFromReader(r.Body, pinUpdate)
	defer r.Body.Close()
	if err != nil {
		logger.Info(err.Error())
		err = responseError(w, "parse_body", "could not read the data to change")
		if err != nil {
			logger.Error(err.Error())
		}
	}
	err = h.pinCase.EditPinByID(r.Context(), int(pinID), userID, pinUpdate)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "edit_pin", "internal error")
	} else {
		err = responseOk(http.StatusOK, w, "pin data has been successfully changed", nil)
	}

	if err != nil {
		logger.Error(err.Error())
	}
}

func (h *HandlerHTTP) ViewPin(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	pinIdStr := chi.URLParam(r, "pinID")
	pinID, err := strconv.ParseInt(pinIdStr, 10, 64)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "parse_url", "internal error")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	userID, ok := r.Context().Value(auth.KeyCurrentUserID).(int)
	if !ok {
		userID = user.UserUnknown
	}
	pin, err := h.pinCase.ViewAnPin(r.Context(), int(pinID), userID)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "edit_pin", "internal error")
	} else {
		err = responseOk(http.StatusOK, w, "pin was successfully received", h.converter.ToPinFromService(pin))
	}
	if err != nil {
		logger.Error(err.Error())
	}
}
