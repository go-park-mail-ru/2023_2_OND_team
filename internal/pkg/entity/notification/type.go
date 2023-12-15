package notification

import "errors"

var ErrUnknownNotifyType = errors.New("unknown notify type")

func TypeString(t NotifyType) string {
	switch t {
	case NotifyComment:
		return "comment"
	case _notifyCustom:
		return "custom"
	}

	return ""
}

func NotifyTemplateByType(t NotifyType) string {
	return notifyTypeTemplate[t]
}
