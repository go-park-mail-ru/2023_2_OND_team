package websocket

import "github.com/go-park-mail-ru/2023_2_OND_team/internal/pkg/entity/message"

const (
	_topicChat         = "chat"
	_topicNotification = "notification"
)

type Channel struct {
	Name string `json:"name"`
}

type Object struct {
	Type    string          `json:"eventType,omitempty"`
	Message message.Message `json:"message"`
}

type PublsihRequest struct {
	ID      int `json:"requestID"`
	Channel Channel
	Message Object
}

type MessageFromChannel struct {
	Type    string          `json:"type"`
	Message ResponseMessage `json:"message"`
}

type ResponseMessage struct {
	Object
	Status      string `json:"status"`
	Code        string `json:"code,omitempty"`
	MessageText string `json:"messageText,omitempty"`
}

type ResponseOnRequest struct {
	ID      int    `json:"requestID"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
	Body    any    `json:"body,omitempty"`
}

func newResponseOnRequest(id int, status, code, message string, body any) *ResponseOnRequest {
	return &ResponseOnRequest{
		ID:      id,
		Type:    "response",
		Status:  status,
		Code:    code,
		Message: message,
		Body:    body,
	}
}

func newMessageFromChannel(status, code string, v any) *MessageFromChannel {
	mes := &MessageFromChannel{
		Type: "event",
		Message: ResponseMessage{
			Status: status,
			Code:   code,
		},
	}
	if v, ok := v.(Object); ok {
		mes.Message.Object = v
		return mes
	}
	if v, ok := v.(string); ok {
		mes.Message.MessageText = v
	}
	return mes
}
