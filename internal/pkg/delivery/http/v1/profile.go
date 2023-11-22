package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	userEntity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/user"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

type UserInfo struct {
	ID           int    `json:"id" example:"123"`
	Username     string `json:"username" example:"Snapshot"`
	Avatar       string `json:"avatar" example:"/pic1"`
	Name         string `json:"name" example:"Bob"`
	Surname      string `json:"surname" example:"Dylan"`
	About        string `json:"about" example:"Cool guy"`
	IsSubscribed bool   `json:"is_subscribed" example:"true"`
	SubsCount    int    `json:"subscribers" example:"23"`
}

type ProfileInfo struct {
	ID        int    `json:"id" example:"1"`
	Username  string `json:"username" example:"baobab"`
	Avatar    string `json:"avatar" example:"/pic1"`
	SubsCount int    `json:"subscribers" example:"12"`
}

func ToUserInfoFromService(user *userEntity.User, isSubscribed bool, subsCount int) UserInfo {
	return UserInfo{
		ID:           user.ID,
		Username:     user.Username,
		Avatar:       user.Avatar,
		Name:         user.Name.String,
		Surname:      user.Surname.String,
		About:        user.AboutMe.String,
		IsSubscribed: isSubscribed,
		SubsCount:    subsCount,
	}
}

func ToProfileInfoFromService(user *userEntity.User, subsCount int) ProfileInfo {
	return ProfileInfo{
		ID:        user.ID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		SubsCount: subsCount,
	}
}

func (h *HandlerHTTP) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userIdParam := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		h.responseErr(w, r, &ErrInvalidUrlParams{map[string]string{"userID": userIdParam}})
		return
	}

	if user, isSubscribed, subsCount, err := h.userCase.GetUserInfo(r.Context(), int(userID)); err != nil {
		h.responseErr(w, r, err)
	} else if err := responseOk(http.StatusOK, w, "got user info successfully", ToUserInfoFromService(user, isSubscribed, subsCount)); err != nil {
		h.responseErr(w, r, err)
	}
}

func (h *HandlerHTTP) GetProfileHeaderInfo(w http.ResponseWriter, r *http.Request) {
	if user, subsCount, err := h.userCase.GetProfileInfo(r.Context()); err != nil {
		h.responseErr(w, r, err)
	} else if err := responseOk(http.StatusOK, w, "got profile info successfully", ToProfileInfoFromService(user, subsCount)); err != nil {
		h.responseErr(w, r, err)
	}
}

func (h *HandlerHTTP) ProfileEditInfo(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	userID := r.Context().Value(auth.KeyCurrentUserID).(int)

	data := &user.ProfileUpdateData{}
	err := json.NewDecoder(r.Body).Decode(data)
	defer r.Body.Close()
	if err != nil {
		logger.Info("json decode: " + err.Error())
		err = responseError(w, "parse_body",
			"the request body must contain json with any of the fields: username, email, name, surname, password")
		if err != nil {
			logger.Error(err.Error())
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
	if data.AboutMe != nil && !isValidAboutMe(*data.AboutMe) {
		invalidFields.addInvalidField("about_me")
	}
	if data.Password != nil && !isValidPassword(*data.Password) {
		invalidFields.addInvalidField("password")
	}
	if invalidFields.Err() != nil {
		err = responseError(w, "invalid_params", invalidFields.Error())
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	err = h.userCase.EditProfileInfo(r.Context(), userID, data)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "uniq_fields", "there is already an account with this username or email")
	} else {
		err = responseOk(http.StatusOK, w, "user data has been successfully changed", nil)
	}

	if err != nil {
		logger.Error(err.Error())
	}
}

func (h *HandlerHTTP) ProfileEditAvatar(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	userID := r.Context().Value(auth.KeyCurrentUserID).(int)
	logger.Info("request on signup", log.F{"method", r.Method}, log.F{"path", r.URL.Path},
		log.F{"userID", fmt.Sprint(userID)}, log.F{"content-type", r.Header.Get("Content-Type")})

	defer r.Body.Close()

	err := h.userCase.UpdateUserAvatar(r.Context(), userID, r.Header.Get("Content-Type"), r.ContentLength, r.Body)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "edit_avatar", "failed to change user's avatar")
	} else {
		err = responseOk(http.StatusOK, w, "the user's avatar has been successfully changed", nil)
	}

	if err != nil {
		logger.Error(err.Error())
	}
}

func (h *HandlerHTTP) GetProfileInfo(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	userID := r.Context().Value(auth.KeyCurrentUserID).(int)
	user, err := h.userCase.GetAllProfileInfo(r.Context(), userID)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "get_info", "failed to get user information")
	} else {
		err = responseOk(http.StatusOK, w, "user data has been successfully received", user)
	}

	if err != nil {
		logger.Error(err.Error())
	}
}
