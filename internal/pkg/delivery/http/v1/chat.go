package v1

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/message"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
)

func (h *HandlerHTTP) FeedChats(w http.ResponseWriter, r *http.Request) {
	log := h.getRequestLogger(r)
	userID := r.Context().Value(auth.KeyCurrentUserID).(int)

	count, lastID, err := FetchValidParamForLoadFeed(r.URL)
	if err != nil {
		log.Info(err.Error())
		err = responseError(w, "parse_url", "bad request url for getting feed chat")
		if err != nil {
			log.Error(err.Error())
		}
		return
	}

	chats, newLastID, err := h.messageCase.GetUserChatsWithOtherUsers(r.Context(), userID, count, lastID)
	if err != nil {
		log.Errorf(err.Error())
	}
	err = responseOk(http.StatusOK, w, "success get feed user chats", map[string]any{
		"chats":  chats,
		"lastID": newLastID,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (h *HandlerHTTP) SendMessageToUser(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)
	userID := r.Context().Value(auth.KeyCurrentUserID).(int)

	fromUserID := r.Context().Value(auth.KeyCurrentUserID).(int)
	toUserID, err := fetchURLParamInt(r, "userID")
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "parse_url", "could not extract to whom the message is being sent")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	mes := &message.Message{}
	err = decodeBody(r, mes)
	defer r.Body.Close()
	if err != nil {
		logger.Info(err.Error())
		err = responseError(w, "parse_body", "invalid request body")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}
	mes.From = fromUserID
	mes.To = toUserID

	idNewMessage, err := h.messageCase.SendMessage(r.Context(), userID, mes)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "send_message", "failed to send message")
	} else {
		err = responseOk(http.StatusCreated, w, "the message was sent successfully",
			map[string]int{"id": idNewMessage})
	}
	if err != nil {
		logger.Error(err.Error())
	}
}

func (h *HandlerHTTP) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	userID := r.Context().Value(auth.KeyCurrentUserID).(int)
	messageID, err := fetchURLParamInt(r, "messageID")
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "parse_url", "could not extract to whom the message is being sent")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	err = h.messageCase.DeleteMessage(r.Context(), userID, &message.Message{ID: messageID})
	if err != nil {
		logger.Warn(err.Error())
		err = responseError(w, "delete_message", "fail deleting a message")
	} else {
		err = responseOk(http.StatusOK, w, "the message was successfully deleted", nil)
	}
	if err != nil {
		logger.Error(err.Error())
	}
}

func (h *HandlerHTTP) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	userID := r.Context().Value(auth.KeyCurrentUserID).(int)
	messageID, err := fetchURLParamInt(r, "messageID")
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "parse_url", "could not extract to whom the message is being sent")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	mes := &message.Message{}
	err = decodeBody(r, mes)
	defer r.Body.Close()
	if err != nil {
		logger.Info(err.Error())
		err = responseError(w, "parse_body", "invalid request body")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}
	mes.ID = messageID

	err = h.messageCase.UpdateContentMessage(r.Context(), userID, mes)
	if err != nil {
		logger.Warn(err.Error())
		err = responseError(w, "update_message", "fail updating a message")
	} else {
		err = responseOk(http.StatusOK, w, "the message was successfully updated", nil)
	}
	if err != nil {
		logger.Error(err.Error())
	}
}

func (h *HandlerHTTP) GetMessagesFromChat(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

	userID := r.Context().Value(auth.KeyCurrentUserID).(int)
	userID2, err := fetchURLParamInt(r, "userID")
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "parse_url", "could not extract to whom the message is being sent")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	count, lastID, err := FetchValidParamForLoadFeed(r.URL)
	if err != nil {
		logger.Error(err.Error())
		err = responseError(w, "bad_request", "failed to get parameters for receiving a message from a chat")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	feed, newLastID, err := h.messageCase.GetMessagesFromChat(r.Context(), userID, message.Chat{userID, userID2}, count, lastID)
	if err != nil {
		logger.Warn(err.Error())
	}
	err = responseOk(http.StatusOK, w, "messages received successfully", map[string]any{
		"messages": h.converter.ToMessagesFromService(feed),
		"lastID":   newLastID,
	})
	if err != nil {
		logger.Error(err.Error())
	}
}
