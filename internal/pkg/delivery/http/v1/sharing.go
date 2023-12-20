package v1

import (
	"net/http"

	"github.com/mailru/easyjson"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/share"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	log "github.com/go-park-mail-ru/2023_2_OND_team/pkg/logger"
)

func (h *HandlerHTTP) CreateSharedLink(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	userID := r.Context().Value(auth.KeyCurrentUserID).(int)
	boardID, err := fetchURLParamInt(r, "boardID")
	if err != nil {
		logger.Info(err.Error())
		err = responseError(w, "parse_url", "the request url could not be get board id")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	sharedLink := &share.SharedLink{}
	err = easyjson.UnmarshalFromReader(r.Body, sharedLink)
	defer r.Body.Close()
	if err != nil {
		logger.Info(err.Error())
		err = responseError(w, "parse_body", "the request body failed parse")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}
	sharedLink.BoardID = boardID

	id, err := h.shareCase.CreateSharedLinkForAddingContributor(r.Context(), userID, sharedLink)
	if err != nil {
		logger.Info(err.Error())
		err = responseError(w, "add_link", "failed to create a link to add contributors to the board with the specified role")
	} else {
		err = responseOk(http.StatusCreated, w, "successful link creation", map[string]int{"id": id})
	}
	if err != nil {
		logger.Error(err.Error())
	}
}

func (h *HandlerHTTP) CheckLink(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)
	userID := r.Context().Value(auth.KeyCurrentUserID).(int)

	boardID, err := fetchURLParamInt(r, "boardID")
	if err != nil {
		logger.Info(err.Error())
		err = responseError(w, "parse_url", "the request url could not be get board id")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	linkID, err := fetchURLParamInt(r, "linkID")
	if err != nil {
		logger.Info(err.Error())
		err = responseError(w, "parse_url", "the request url could not be get link id")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	_, err = h.shareCase.CheckLinkAvailability(r.Context(), userID, linkID)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "check_link", "the link is unavailable to the user")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	board, username, err := h.boardCase.GetBoardWithAuthor(r.Context(), boardID)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "link_info", "failed to get information about the link")
	} else {
		err = responseOk(http.StatusOK, w, "", map[string]any{"board": board, "author": username})
	}

	if err != nil {
		logger.Error(err.Error())
	}
}

func (h *HandlerHTTP) SharedBoard(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)
	userID := r.Context().Value(auth.KeyCurrentUserID).(int)

	boardID, err := fetchURLParamInt(r, "boardID")
	if err != nil {
		logger.Info(err.Error())
		err = responseError(w, "parse_url", "the request url could not be get board id")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	linkID, err := fetchURLParamInt(r, "linkID")
	if err != nil {
		logger.Info(err.Error())
		err = responseError(w, "parse_url", "the request url could not be get link id")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	role, err := h.shareCase.CheckLinkAvailability(r.Context(), userID, linkID)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "check_link", "the link is unavailable to the user")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	logger.Info("get link role", log.F{"role", role})

	err = h.boardCase.AddContributorsToBoard(r.Context(), boardID, []int{userID}, role)
	if err != nil {
		logger.Info(err.Error())
		err = responseError(w, "add_contributor", "the user could not be added to the board contributors")
	} else {
		err = responseOk(http.StatusCreated, w, "the user has successfully become a contributor to the board", nil)
	}

	if err != nil {
		logger.Error(err.Error())
	}
}
