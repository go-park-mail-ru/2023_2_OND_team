package v1

import (
	"net/http"
	"strconv"

	errHTTP "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1/errors"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/delivery/http/v1/structs"
	userEntity "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	"github.com/mailru/easyjson"
)

var (
	defaultSubCount   = 20
	defaultSubLastID  = 1 << 30
	subscriptionsView = "subscriptions"
	subscribersView   = "subscribers"
	maxCount          = 50
)

func (h *HandlerHTTP) Subscribe(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != ApplicationJson {
		h.responseErr(w, r, &errHTTP.ErrInvalidContentType{PreferredType: ApplicationJson})
		return
	}

	sub := structs.SubscriptionAction{}

	if err := easyjson.UnmarshalFromReader(r.Body, &sub); err != nil {
		h.responseErr(w, r, &errHTTP.ErrInvalidBody{})
		return
	}
	defer r.Body.Close()
	if err := sub.Validate(); err != nil {
		h.responseErr(w, r, err)
		return
	}

	from := r.Context().Value(auth.KeyCurrentUserID).(int)
	if err := h.subCase.SubscribeToUser(r.Context(), from, *sub.To); err != nil {
		h.responseErr(w, r, err)
	} else if err := responseOk(http.StatusOK, w, "subscribed successfully", nil); err != nil {
		h.responseErr(w, r, err)
	}

}

func (h *HandlerHTTP) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != ApplicationJson {
		h.responseErr(w, r, &errHTTP.ErrInvalidContentType{PreferredType: ApplicationJson})
		return
	}

	sub := structs.SubscriptionAction{}
	if err := easyjson.UnmarshalFromReader(r.Body, &sub); err != nil {
		h.responseErr(w, r, &errHTTP.ErrInvalidBody{})
		return
	}
	defer r.Body.Close()
	if err := sub.Validate(); err != nil {
		h.responseErr(w, r, err)
		return
	}

	from := r.Context().Value(auth.KeyCurrentUserID).(int)
	if err := h.subCase.UnsubscribeFromUser(r.Context(), from, *sub.To); err != nil {
		h.responseErr(w, r, err)
	} else if err := responseOk(http.StatusOK, w, "unsubscribed successfully", nil); err != nil {
		h.responseErr(w, r, err)
	}
}

func (h *HandlerHTTP) GetSubscriptionInfoForUser(w http.ResponseWriter, r *http.Request) {
	opts, err := GetOpts(r)
	if err != nil {
		h.responseErr(w, r, err)
		return
	}

	if users, err := h.subCase.GetSubscriptionInfoForUser(r.Context(), opts); err != nil {
		h.responseErr(w, r, err)
	} else if err := responseOk(http.StatusOK, w, "got subscription info successfully", h.converter.ToSubscriptionUsersFromService(users)); err != nil {
		h.responseErr(w, r, err)
	}
}

func GetOpts(r *http.Request) (*userEntity.SubscriptionOpts, error) {
	opts := &userEntity.SubscriptionOpts{}
	invalidParams := map[string]string{}

	var (
		userID, count, lastID int64
		filter                string
		err                   error
	)
	if userIdParam := r.URL.Query().Get("userID"); userIdParam != "" {
		if userID, err = strconv.ParseInt(userIdParam, 10, 64); err != nil || userID < 0 {
			invalidParams["userID"] = userIdParam
		} else {
			opts.UserID = int(userID)
		}
	} else {
		opts.UserID, _ = r.Context().Value(auth.KeyCurrentUserID).(int)
	}

	if countParam := r.URL.Query().Get("count"); countParam != "" {
		if count, err = strconv.ParseInt(countParam, 10, 64); err != nil || count < 0 {
			invalidParams["count"] = countParam
		} else {
			opts.Count = int(count)
		}
	} else {
		opts.Count = defaultSubCount
	}

	if lastIdParam := r.URL.Query().Get("lastID"); lastIdParam != "" {
		if lastID, err = strconv.ParseInt(lastIdParam, 10, 64); err != nil || lastID < 0 {
			invalidParams["lastID"] = lastIdParam
		} else {
			opts.LastID = int(lastID)
		}
	} else {
		opts.LastID = defaultSubLastID
	}

	if filter = r.URL.Query().Get("view"); filter != "" {
		if filter != subscriptionsView && filter != subscribersView {
			invalidParams["view"] = filter
		} else {
			opts.Filter = filter
		}
	} else {
		invalidParams["view"] = filter
	}

	if opts.Count > maxCount {
		opts.Count = maxCount
	}
	if len(invalidParams) > 0 {
		return nil, &errHTTP.ErrInvalidQueryParam{invalidParams}
	}
	return opts, nil
}
