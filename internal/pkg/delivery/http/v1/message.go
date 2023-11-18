package v1

import (
	"net/http"

	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/message"
	"github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/middleware/auth"
)

func (h *HandlerHTTP) SendMessageToUser(w http.ResponseWriter, r *http.Request) {
	logger := h.getRequestLogger(r)

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
}

func (h *HandlerHTTP) DeleteMessage(w http.ResponseWriter, r *http.Request) {

}

func (h *HandlerHTTP) UpdateMessage(w http.ResponseWriter, r *http.Request) {

}

func (h *HandlerHTTP) GetMessagesFromChat(w http.ResponseWriter, r *http.Request) {

}
