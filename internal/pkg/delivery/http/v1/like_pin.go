package v1

import (
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func (h *HandlerHTTP) SetLikePin(w http.ResponseWriter, r *http.Request) {
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

	countLike, err := h.pinCase.SetLikeFromUser(r.Context(), int(pinID), userID)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "like_pin_set", "internal error")
	} else {
		err = responseOk(http.StatusCreated, w, "ok", map[string]int{"count_like": countLike})
	}
	if err != nil {
		logger.Error(err.Error())
	}
}

func (h *HandlerHTTP) DeleteLikePin(w http.ResponseWriter, r *http.Request) {
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

	countLike, err := h.pinCase.DeleteLikeFromUser(r.Context(), int(pinID), userID)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "like_pin_del", "internal error")
	} else {
		err = responseOk(http.StatusOK, w, "ok", map[string]int{"count_like": countLike})
	}
	if err != nil {
		logger.Error(err.Error())
	}
}

func (h *HandlerHTTP) IsSetLikePin(w http.ResponseWriter, r *http.Request) {
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

	isSet, err := h.pinCase.CheckUserHasSetLike(r.Context(), int(pinID), userID)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "like_pin_set", "internal error")
	} else {
		err = responseOk(http.StatusOK, w, "ok", map[string]bool{"is_set": isSet})
	}
	if err != nil {
		logger.Error(err.Error())
	}
}

// unused
func (h *HandlerHTTP) LikedPins(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)
	userID := r.Context().Value(auth.KeyCurrentUserID).(int)

	count, minID, maxID, err := FetchValidParamForLoadTape(r.URL)
	if err != nil {
		logger.Info("parse url query params", log.F{"error", err.Error()})
		err = responseError(w, "bad_params",
			"expected parameters: count(positive integer: [1; 1000]), maxID, minID(positive integers, the absence of these parameters is equal to the value 0)")
	} else {
		logger.Infof("param: count=%d, minID=%d, maxID=%d", count, minID, maxID)
		pins, minID, maxID := h.pinCase.SelectUserPins(r.Context(), userID, count, minID, maxID)
		err = responseOk(http.StatusOK, w, "pins received are sorted by id", map[string]any{
			"pins":  pins,
			"minID": minID,
			"maxID": maxID,
		})
	}
	if err != nil {
		logger.Error(err.Error())
	}

}
