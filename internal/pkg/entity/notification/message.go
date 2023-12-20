package notification

//go:generate easyjson
//easyjson:json
type NotifyMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	err     error
}

func (n *NotifyMessage) Err() error {
	return n.err
}

func NewNotifyMessage(t NotifyType, content string) *NotifyMessage {
	return &NotifyMessage{
		Type:    TypeString(t),
		Content: content,
	}
}

func NewNotifyMessageWithError(err error) *NotifyMessage {
	return &NotifyMessage{
		err: err,
	}
}
