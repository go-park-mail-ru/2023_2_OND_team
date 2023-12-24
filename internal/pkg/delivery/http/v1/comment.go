package v1

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/comment"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/user"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
	"github.com/mailru/easyjson"
)

func (h *HandlerHTTP) WriteComment(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)
	userID := r.Context().Value(auth.KeyCurrentUserID).(int)

	pinID, err := fetchURLParamInt(r, "pinID")
	if err != nil {
		err = responseError(w, "parse_url", "the request url could not be get pin id")
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}

	comment := &comment.Comment{}
	err = easyjson.UnmarshalFromReader(r.Body, comment)
	defer r.Body.Close()
	if err != nil {
		logger.Warn(err.Error())

		err = responseError(w, "parse_body", "the request body could not be parsed to send a comment")
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}

	comment.PinID = pinID
	_, err = h.commentCase.PutCommentOnPin(r.Context(), userID, comment)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "create_comment", "couldn't leave a comment under the selected pin")
	} else {
		err = responseOk(http.StatusCreated, w, "the comment has been added successfully", nil)
	}
	if err != nil {
		logger.Error(err.Error())
	}
}

func (h *HandlerHTTP) DeleteComment(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)
	userID := r.Context().Value(auth.KeyCurrentUserID).(int)

	commentID, err := fetchURLParamInt(r, "commentID")
	if err != nil {
		err = responseError(w, "parse_url", "the request url could not be get pin id")
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}

	err = h.commentCase.DeleteComment(r.Context(), userID, commentID)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "delete_comment", "couldn't delete pin comment")
	} else {
		err = responseOk(http.StatusOK, w, "the comment was successfully deleted", nil)
	}
	if err != nil {
		logger.Error(err.Error())
	}
}

func (h *HandlerHTTP) ViewFeedComment(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)
	userID, ok := r.Context().Value(auth.KeyCurrentUserID).(int)
	if !ok {
		userID = user.UserUnknown
	}

	pinID, err := fetchURLParamInt(r, "pinID")
	if err != nil {
		err = responseError(w, "parse_url", "the request url could not be get pin id")
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}

	count, lastID, err := FetchValidParamForLoadFeed(r.URL)
	if err != nil {
		err = responseError(w, "query_param", "the parameters for displaying the pin feed could not be extracted from the request")
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}

	feed, newLastID, err := h.commentCase.GetFeedCommentOnPin(r.Context(), userID, pinID, count, lastID)
	if err != nil && len(feed) == 0 {
		err = responseError(w, "feed_view", "error displaying pin comments")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	if err != nil {
		logger.Error(err.Error())
	}

	err = responseOk(http.StatusOK, w, "feed comment to pin", map[string]any{"comments": h.converter.ToCommentsFromService(feed), "lastID": newLastID})
	if err != nil {
		logger.Error(err.Error())
	}
}
