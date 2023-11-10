package v1

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	usecase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/pin"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

// count, minID, maxID, liked{true,false}, protection{private,public,all}, board_id, user_id
func (h *HandlerHTTP) FeedPins(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)
	userID, isAuth := r.Context().Value(auth.KeyCurrentUserID).(int)
	if !isAuth {
		userID = usecase.UserUnknown
	}

	logger.Info("request on getting feed of pins", log.F{"rawQuery", r.URL.RawQuery})

	cfg := parseFeedConfig(r.URL)

	feed, err := h.pinCase.ViewFeedPin(r.Context(), userID, cfg)
	logger.Info("send feed pins", log.F{"count", len(feed.Pins)})

	if err != nil {
		err = responseError(w, "fsdf", "dsfsdf")
	} else {
		err = responseOk(http.StatusOK, w, "ok", feed)
	}
	if err != nil {
		logger.Error(err.Error())
	}
}

func parseFeedConfig(u *url.URL) pin.FeedPinConfig {
	cfg := pin.FeedPinConfig{}
	i, _ := strconv.ParseInt(u.Query().Get("minID"), 10, 64)
	cfg.MinID = int(i)

	i, _ = strconv.ParseInt(u.Query().Get("maxID"), 10, 64)
	cfg.MaxID = int(i)

	i, _ = strconv.ParseInt(u.Query().Get("userID"), 10, 64)
	cfg.UserID = int(i)

	i, _ = strconv.ParseInt(u.Query().Get("boardID"), 10, 64)
	cfg.BoardID = int(i)

	i, _ = strconv.ParseInt(u.Query().Get("count"), 10, 64)
	cfg.Count = int(i)

	ok, _ := strconv.ParseBool(u.Query().Get("deleted"))
	cfg.Deleted = ok

	ok, _ = strconv.ParseBool(u.Query().Get("liked"))
	cfg.Liked = ok

	switch u.Query().Get("protection") {
	case "public":
		cfg.Protection = pin.FeedProtectionPublic
	case "private":
		cfg.Protection = pin.FeedProtectionPrivate
	default:
		cfg.Protection = pin.FeedAll
	}

	return cfg
}
