package v1

import (
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func (h *HandlerHTTP) SetLikePin(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request on set like for pin", log.F{"method", r.Method}, log.F{"path", r.URL.Path})

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

	countLike, err := h.pinCase.SetLikeFromUser(r.Context(), int(pinID), userID)
	if err != nil {
		h.log.Error(err.Error())
		err = responseError(w, "like_pin_set", "internal error")
	} else {
		err = responseOk(w, "ok", map[string]int{"count_like": countLike})
	}
	if err != nil {
		h.log.Error(err.Error())
	}
}

func (h *HandlerHTTP) DeleteLikePin(w http.ResponseWriter, r *http.Request) {
	h.log.Info("request on set like for pin", log.F{"method", r.Method}, log.F{"path", r.URL.Path})

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

	err = h.pinCase.DeleteLikeFromUser(r.Context(), int(pinID), userID)
	if err != nil {
		h.log.Error(err.Error())
		err = responseError(w, "like_pin_del", "internal error")
	} else {
		err = responseOk(w, "ok", nil)
	}
	if err != nil {
		h.log.Error(err.Error())
	}
}
