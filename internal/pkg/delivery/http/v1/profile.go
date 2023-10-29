package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func (h *HandlerHTTP) ProfileEditInfo(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request on signup", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	userID := r.Context().Value(auth.KeyCurrentUserID).(int)

	data := user.NewProfileUpdateData()
	err := json.NewDecoder(r.Body).Decode(data)
	defer r.Body.Close()
	if err != nil {
		h.log.Info("json decode: " + err.Error())
		err = responseError(w, "parse_body",
			"the request body must contain json with any of the fields: username, email, name, surname, password")
		if err != nil {
			h.log.Error(err.Error())
		}
		return
	}

	invalidFields := new(errorFields)
	if data.Username != nil && !isValidUsername(*data.Username) {
		invalidFields.addInvalidField("username")
	}
	if data.Email != nil && !isValidEmail(*data.Email) {
		invalidFields.addInvalidField("email")
	}
	if data.Name != nil && !isValidName(*data.Name) {
		invalidFields.addInvalidField("name")
	}
	if data.Surname != nil && !isValidSurname(*data.Surname) {
		invalidFields.addInvalidField("surname")
	}
	if data.Password != nil && !isValidPassword(*data.Password) {
		invalidFields.addInvalidField("password")
	}
	if invalidFields.Err() != nil {
		err = responseError(w, "invalid_params", err.Error())
		if err != nil {
			h.log.Error(err.Error())
		}
		return
	}

	err = h.userCase.EditProfileInfo(r.Context(), userID, data)
	if err != nil {
		h.log.Error(err.Error())
		err = responseError(w, "uniq_fields", "there is already an account with this username or email")
	} else {
		err = responseOk(w, "user data has been successfully changed", nil)
	}

	if err != nil {
		h.log.Error(err.Error())
	}
}

func (h *HandlerHTTP) ProfileEditAvatar(w http.ResponseWriter, r *http.Request) {
	SetContentTypeJSON(w)

	userID := r.Context().Value(auth.KeyCurrentUserID).(int)
	h.log.Info("request on signup", log.F{"method", r.Method}, log.F{"path", r.URL.Path},
		log.F{"userID", fmt.Sprint(userID)}, log.F{"content-type", r.Header.Get("Content-Type")})

	defer r.Body.Close()

	err := h.userCase.UpdateUserAvatar(r.Context(), userID, r.Body, r.Header.Get("Content-Type"))
	if err != nil {
		h.log.Error(err.Error())
		err = responseError(w, "edit_avatar", "failed to change user's avatar")
	} else {
		err = responseOk(w, "the user's avatar has been successfully changed", nil)
	}

	if err != nil {
		h.log.Error(err.Error())
	}
}

func (h *HandlerHTTP) GetProfileInfo(w http.ResponseWriter, r *http.Request) {
	SetContentTypeJSON(w)

	userID := r.Context().Value(auth.KeyCurrentUserID).(int)
	user, err := h.userCase.GetAllProfileInfo(r.Context(), userID)
	if err != nil {
		h.log.Error(err.Error())
		err = responseError(w, "get_info", "failed to get user information")
	} else {
		err = responseOk(w, "user data has been successfully received", user)
	}

	if err != nil {
		h.log.Error(err.Error())
	}
}