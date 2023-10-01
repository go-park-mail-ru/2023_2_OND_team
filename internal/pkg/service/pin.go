package service

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

var ErrCountParameterMissing = errors.New("the count parameter is missing")
var ErrBadParams = errors.New("bad params")

func (s *Service) GetPins(w http.ResponseWriter, r *http.Request) {
	s.log.Info("request on get pins", log.F{"method", r.Method}, log.F{"path", r.URL.Path})
	SetContentTypeJSON(w)

	count, lastID, err := fetchValidParamForLoadTape(r.URL)
	if err != nil {
		s.log.Info("parse body error", log.F{"error", err.Error()})
		fmt.Fprintln(w, "{\"status\": \"error\"}")
		return
	}

	s.log.Sugar().Infof("param: count=%d, lastID=%d", count, lastID)

	pins, last := s.pinCase.SelectNewPins(r.Context(), count, lastID)
	fmt.Fprintf(w, `{"status": "ok",
	"message": "download new pins",
	"body": {
		"pins": %v,
		"lastID": %d
	}
	}`, pins, last)
}

func fetchValidParamForLoadTape(u *url.URL) (count int, lastID int, err error) {
	if param := u.Query().Get("count"); len(param) > 0 {
		c, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("fetch count param for load tape: %w", err)
		}
		count = int(c)
	} else {
		return 0, 0, ErrCountParameterMissing
	}
	if param := u.Query().Get("lastID"); len(param) > 0 {
		last, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("fetch lastID param for load tape: %w", err)
		}
		lastID = int(last)
	}
	if count <= 0 || lastID < 0 {
		return 0, 0, ErrBadParams
	}
	return
}
