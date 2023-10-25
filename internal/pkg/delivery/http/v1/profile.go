package v1

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func (h *HandlerHTTP) ProfileEditInfo(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request on signup", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
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
		responseError(w, "edit_avatar", "failed to change user's avatar")
	} else {
		responseOk(w, "the user's avatar has been successfully changed", nil)
	}
}

func (h *HandlerHTTP) GetProfileInfo(w http.ResponseWriter, r *http.Request) {
	SetContentTypeJSON(w)

	userID := r.Context().Value(auth.KeyCurrentUserID).(int)
	user, err := h.userCase.GetAllProfileInfo(r.Context(), userID)
	if err != nil {
		h.log.Error(err.Error())
		responseError(w, "get_info", "failed to get user information")
	} else {
		responseOk(w, "user data has been successfully received", user)
	}
}
