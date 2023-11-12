package v1

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/pin"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	usecase "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/usecase/pin"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func (h *HandlerHTTP) FeedPins(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)
	userID, isAuth := r.Context().Value(auth.KeyCurrentUserID).(int)
	if !isAuth {
		userID = usecase.UserUnknown
	}

	logger.Info("request on getting feed of pins", log.F{"rawQuery", r.URL.RawQuery})

	cfg, err := parseFeedConfig(r.URL)
	if err != nil {
		logger.Info("error parse query params", log.F{"parse_error", err.Error()})
		err = responseError(w, "parse_params", "bad url params")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	err = h.boardCase.CheckAvailabilityFeedPinCfgOnBoard(r.Context(), cfg, userID, isAuth)
	if err != nil {
		logger.Info(err.Error())
		err = responseError(w, "no_access", "there is no access to get board pins")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}
	feed, err := h.pinCase.ViewFeedPin(r.Context(), userID, cfg)
	logger.Info("send feed pins", log.F{"count", len(feed.Pins)})

	if err != nil {
		err = responseError(w, "no_access", "there is no access to get board pins")
	} else {
		err = responseOk(http.StatusOK, w, "ok", feed)
	}
	if err != nil {
		logger.Error(err.Error())
	}
}

func parseFeedConfig(u *url.URL) (pin.FeedPinConfig, error) {
	cfg := pin.FeedPinConfig{}
	var (
		err      error
		numInt64 int64
		ok       bool
	)

	if !u.Query().Has("count") {
		return cfg, errors.New("parse feed config: require count")
	}
	numInt64, err = strconv.ParseInt(u.Query().Get("count"), 10, 64)
	if err != nil {
		return cfg, fmt.Errorf("pars feed config: %w", err)
	}
	cfg.Count = int(numInt64)

	if u.Query().Has("minID") {
		numInt64, err = strconv.ParseInt(u.Query().Get("minID"), 10, 64)
		if err != nil {
			return pin.FeedPinConfig{}, fmt.Errorf("pars feed config: %w", err)
		}
		cfg.MinID = int(numInt64)
	}

	if u.Query().Has("maxID") {
		numInt64, err = strconv.ParseInt(u.Query().Get("maxID"), 10, 64)
		if err != nil {
			return pin.FeedPinConfig{}, fmt.Errorf("pars feed config: %w", err)
		}
		cfg.MaxID = int(numInt64)
	}

	if u.Query().Has("userID") {
		numInt64, err = strconv.ParseInt(u.Query().Get("userID"), 10, 64)
		if err != nil {
			return pin.FeedPinConfig{}, fmt.Errorf("pars feed config: %w", err)
		}
		cfg.SetUser(int(numInt64))
	}

	if u.Query().Has("boardID") {
		numInt64, err = strconv.ParseInt(u.Query().Get("boardID"), 10, 64)
		if err != nil {
			return pin.FeedPinConfig{}, fmt.Errorf("pars feed config: %w", err)
		}
		cfg.SetBoard(int(numInt64))
	}

	if u.Query().Has("deleted") {
		ok, err = strconv.ParseBool(u.Query().Get("deleted"))
		if err != nil {
			return cfg, fmt.Errorf("pars feed config: %w", err)
		}
		cfg.Deleted = ok
	}

	if u.Query().Has("liked") {
		ok, err = strconv.ParseBool(u.Query().Get("liked"))
		if err != nil {
			return cfg, fmt.Errorf("pars feed config: %w", err)
		}
		cfg.Liked = ok
	}

	switch u.Query().Get("protection") {
	case "all":
		cfg.Protection = pin.FeedAll
	case "private":
		cfg.Protection = pin.FeedProtectionPrivate
	default:
		cfg.Protection = pin.FeedProtectionPublic
	}

	return cfg, nil
}
